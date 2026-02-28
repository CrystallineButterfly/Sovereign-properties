import React, { useMemo, useState, useEffect, useRef } from "react";
import { useAuth } from "./AuthProvider";
import { useUXMode } from "./UXModeProvider";
import { apiClient } from "../../../shared/src/utils/api";
import {
  Client as XMTPClient,
  IdentifierKind,
} from "@xmtp/browser-sdk";
import type {
  House,
  Bill,
  Listing,
  RentRequestPayload,
  ConversationSummary,
  ConversationDetails,
} from "../../../shared/src/types";
import toast from "react-hot-toast";
import { Link, useSearchParams } from "react-router-dom";
import { ethers } from "ethers";
import { saveLatestClaimKeyHash } from "../utils/claimKeyStorage";

interface HouseDetailsProps {
  tokenId: string;
}

type WorkflowStepStatus = "completed" | "current" | "upcoming";
type DetailTabKey = "details" | "bills" | "rental" | "messages";

interface WorkflowQuickAction {
  label: string;
  to?: string;
  tab?: DetailTabKey;
  primary?: boolean;
}

const WORKFLOW_STEP_STYLES: Record<
  WorkflowStepStatus,
  {
    label: string;
    card: string;
    badge: string;
    index: string;
  }
> = {
  completed: {
    label: "Completed",
    card: "border-emerald-400/45 bg-emerald-500/10",
    badge: "border-emerald-300/45 bg-emerald-500/20 text-emerald-200",
    index: "border-emerald-300/50 bg-emerald-500/20 text-emerald-100",
  },
  current: {
    label: "Current",
    card: "border-sky-400/50 bg-sky-500/10 shadow-[0_0_24px_rgba(56,189,248,0.16)]",
    badge: "border-sky-300/50 bg-sky-500/20 text-sky-200",
    index: "border-sky-300/60 bg-sky-500/20 text-sky-100",
  },
  upcoming: {
    label: "Upcoming",
    card: "border-slate-600/50 bg-slate-900/35",
    badge: "border-slate-500/50 bg-slate-800/50 text-slate-300",
    index: "border-slate-500/50 bg-slate-800/50 text-slate-300",
  },
};

const parseChainIdValue = (value: string | null): number | null => {
  if (!value) return null;
  const trimmed = value.trim();
  if (!trimmed) return null;
  const parts = trimmed.split(":");
  const parsed = Number.parseInt(parts[parts.length - 1], 10);
  return Number.isFinite(parsed) ? parsed : null;
};

const formatWeiAsEth = (wei: string): string => {
  try {
    return `${ethers.formatEther(wei)} ETH`;
  } catch {
    return `${wei} wei`;
  }
};

type XMTPEnv = "local" | "dev" | "production";

type XmtpClientInstance = Awaited<ReturnType<typeof XMTPClient.create>>;

type XmtpDmConversation = Awaited<
  ReturnType<XmtpClientInstance["conversations"]["createDmWithIdentifier"]>
>;

export const HouseDetails: React.FC<HouseDetailsProps> = ({ tokenId }) => {
  const { walletAddress, chainId, getEthereumProvider } = useAuth();
  const { mode } = useUXMode();
  const [searchParams] = useSearchParams();
  const [house, setHouse] = useState<House | null>(null);
  const [bills, setBills] = useState<Bill[]>([]);
  const [listing, setListing] = useState<Listing | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<DetailTabKey>("details");
  const [actionLoading, setActionLoading] = useState<
    "cancel" | "buy" | "rent" | null
  >(null);
  const [rentDurationDays, setRentDurationDays] = useState("30");
  const [rentDepositAmount, setRentDepositAmount] = useState("");
  const [conversations, setConversations] = useState<ConversationSummary[]>([]);
  const [activeConversationId, setActiveConversationId] = useState("");
  const [conversationMessages, setConversationMessages] = useState<
    ConversationDetails["messages"]
  >([]);
  const [selectedRecipient, setSelectedRecipient] = useState("");
  const [draftMessage, setDraftMessage] = useState("");
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [isSendingMessage, setIsSendingMessage] = useState(false);
  const [xmtpStatus, setXmtpStatus] = useState<
    "idle" | "connecting" | "ready" | "error"
  >("idle");
  const xmtpClientRef = useRef<XmtpClientInstance | null>(null);
  const xmtpWalletAddressRef = useRef("");
  const xmtpReachabilityCacheRef = useRef(new Map<string, boolean>());

  const expectedChainId = Number.parseInt(
    String(import.meta.env.VITE_EXPECTED_CHAIN_ID || ""),
    10,
  );
  const connectedChainId = parseChainIdValue(chainId);
  const wrongChain =
    Number.isFinite(expectedChainId) &&
    connectedChainId !== null &&
    connectedChainId !== expectedChainId;

  const notifyWrongNetwork = () => {
    if (!Number.isFinite(expectedChainId)) {
      return;
    }
    toast.error(`Wrong network. Switch to chain ${expectedChainId} first.`);
  };

  const isOwner = useMemo(() => {
    if (!walletAddress || !house?.ownerAddress) return false;
    return house.ownerAddress.toLowerCase() === walletAddress.toLowerCase();
  }, [house?.ownerAddress, walletAddress]);

  const messageRecipients = useMemo(() => {
    if (!walletAddress || !house) {
      return [];
    }

    const normalizedViewer = walletAddress.toLowerCase();
    const recipients = new Map<string, string>();
    const addRecipient = (address: string | undefined, roleLabel: string) => {
      if (!address) {
        return;
      }
      const normalized = address.toLowerCase();
      if (!/^0x[a-fA-F0-9]{40}$/.test(address) || normalized === normalizedViewer) {
        return;
      }
      if (!recipients.has(normalized)) {
        recipients.set(normalized, roleLabel);
      }
    };

    const owner = house.ownerAddress;
    const originalOwner = house.originalOwner;
    const allowedBuyer = house.listing?.allowedBuyer;
    const renter = house.rental?.renterAddress;
    const isViewerSeller =
      normalizedViewer === owner.toLowerCase()
      || normalizedViewer === originalOwner.toLowerCase();
    const isViewerBuyer = Boolean(
      allowedBuyer && normalizedViewer === allowedBuyer.toLowerCase(),
    );
    const isViewerRenter = Boolean(
      renter && normalizedViewer === renter.toLowerCase(),
    );

    if (isViewerSeller) {
      addRecipient(allowedBuyer, "Buyer");
      addRecipient(renter, "Renter");
    }
    if (isViewerBuyer) {
      addRecipient(owner, "Seller");
      addRecipient(originalOwner, "Seller");
    }
    if (isViewerRenter) {
      addRecipient(owner, "Landlord");
    }
    const listingType = house.listing?.listingType || "none";
    const isViewerNonOwner =
      normalizedViewer !== owner.toLowerCase() &&
      normalizedViewer !== originalOwner.toLowerCase();
    if (isViewerNonOwner && listingType === "for_sale") {
      addRecipient(owner, "Seller");
    }
    if (isViewerNonOwner && listingType === "for_rent") {
      addRecipient(owner, "Landlord");
    }
    for (const conversation of conversations) {
      const counterpartWallet = String(
        conversation.counterpartWalletAddress || "",
      ).trim();
      const roleLabel = conversation.counterpartRole
        ? `${conversation.counterpartRole.charAt(0).toUpperCase()}${conversation.counterpartRole.slice(1)}`
        : "Counterparty";
      addRecipient(counterpartWallet, roleLabel);
    }

    return Array.from(recipients.entries()).map(([address, role]) => ({
      address,
      role,
    }));
  }, [conversations, house, walletAddress]);

  useEffect(() => {
    const requestedTab = searchParams.get("tab");
    if (
      requestedTab === "details" ||
      requestedTab === "bills" ||
      requestedTab === "rental" ||
      requestedTab === "messages"
    ) {
      setActiveTab(requestedTab);
    }

    const requestedRecipient = String(searchParams.get("to") || "")
      .trim()
      .toLowerCase();
    if (/^0x[a-fA-F0-9]{40}$/.test(requestedRecipient)) {
      setSelectedRecipient(requestedRecipient);
    }

    const requestedConversation = String(
      searchParams.get("conversation") || "",
    ).trim();
    if (requestedConversation) {
      setActiveConversationId(requestedConversation);
    }
  }, [searchParams]);

  const resolveXmtpEnv = (): XMTPEnv => {
    const configured = String(import.meta.env.VITE_XMTP_ENV || "")
      .trim()
      .toLowerCase();
    if (configured === "local" || configured === "dev" || configured === "production") {
      return configured;
    }
    const isLocalhost =
      typeof window !== "undefined"
      && (window.location.hostname === "localhost"
        || window.location.hostname === "127.0.0.1");
    return isLocalhost ? "dev" : "production";
  };

  const getXmtpClient = async (): Promise<XmtpClientInstance> => {
    const normalizedWalletAddress = String(walletAddress || "").toLowerCase();
    if (!normalizedWalletAddress || !/^0x[a-fA-F0-9]{40}$/.test(normalizedWalletAddress)) {
      throw new Error("Connect a valid wallet before using XMTP.");
    }

    if (
      xmtpClientRef.current &&
      xmtpWalletAddressRef.current === normalizedWalletAddress
    ) {
      return xmtpClientRef.current;
    }

    setXmtpStatus("connecting");
    const ethereumProvider = await getEthereumProvider();
    const provider = new ethers.BrowserProvider(ethereumProvider);
    const signer = await provider.getSigner();
    const signerAddress = (await signer.getAddress()).toLowerCase();

    if (signerAddress !== normalizedWalletAddress) {
      throw new Error(
        "Active signer wallet does not match the connected application wallet.",
      );
    }

    const xmtpSigner = {
      type: "EOA" as const,
      getIdentifier: () => ({
        identifier: signerAddress,
        identifierKind: IdentifierKind.Ethereum,
      }),
      signMessage: async (message: string) => {
        const signature = await signer.signMessage(message);
        return new Uint8Array(ethers.getBytes(signature));
      },
    };

    const client = await XMTPClient.create(xmtpSigner, {
      env: resolveXmtpEnv(),
    });
    xmtpClientRef.current = client;
    xmtpWalletAddressRef.current = signerAddress;
    setXmtpStatus("ready");
    return client;
  };

  const getXmtpDmConversation = async (
    client: XmtpClientInstance,
    recipientWalletAddress: string,
  ): Promise<XmtpDmConversation> => {
    const normalizedRecipient = recipientWalletAddress.toLowerCase();
    const cachedCanMessage = xmtpReachabilityCacheRef.current.get(
      normalizedRecipient,
    );
    let canMessageRecipient = cachedCanMessage;
    const recipientIdentifier = {
      identifier: normalizedRecipient,
      identifierKind: IdentifierKind.Ethereum,
    };
    if (canMessageRecipient === undefined) {
      const canMessageMap = await client.canMessage([recipientIdentifier]);
      canMessageRecipient = Boolean(
        canMessageMap.get(normalizedRecipient) ??
          canMessageMap.get(recipientIdentifier.identifier),
      );
      xmtpReachabilityCacheRef.current.set(
        normalizedRecipient,
        canMessageRecipient,
      );
    }
    if (!canMessageRecipient) {
      throw new Error(
        `Recipient wallet is not XMTP-enabled on ${resolveXmtpEnv()} yet. ` +
          "Ask them to open messaging once, then retry.",
      );
    }

    const existingConversation = await client.conversations.fetchDmByIdentifier(
      recipientIdentifier,
    );
    if (existingConversation) {
      return existingConversation as XmtpDmConversation;
    }

    return client.conversations.createDmWithIdentifier(recipientIdentifier);
  };

  const resolveInboxWalletMap = async (
    client: XmtpClientInstance,
    inboxIds: string[],
  ): Promise<Map<string, string>> => {
    const map = new Map<string, string>();
    if (inboxIds.length === 0) {
      return map;
    }
    try {
      const states = await client.preferences.fetchInboxStates(inboxIds);
      for (const state of states) {
        const ethereumIdentifier = state.accountIdentifiers.find(
          (identifier) => identifier.identifierKind === IdentifierKind.Ethereum,
        );
        if (ethereumIdentifier?.identifier) {
          map.set(state.inboxId, ethereumIdentifier.identifier.toLowerCase());
        }
      }
      return map;
    } catch {
      try {
        const states = await client.preferences.getInboxStates(inboxIds);
        for (const state of states) {
          const ethereumIdentifier = state.accountIdentifiers.find(
            (identifier) => identifier.identifierKind === IdentifierKind.Ethereum,
          );
          if (ethereumIdentifier?.identifier) {
            map.set(state.inboxId, ethereumIdentifier.identifier.toLowerCase());
          }
        }
      } catch {
        // Ignore identity resolution failures; fallback mapping is used.
      }
      return map;
    }
  };

  const loadMessagesFromXmtp = async (
    recipientWalletAddress: string,
    conversationId: string,
  ) => {
    const client = await getXmtpClient();
    const dmConversation = await getXmtpDmConversation(client, recipientWalletAddress);
    const rawMessages = await dmConversation.messages();

    const inboxIds = Array.from(
      new Set(
        rawMessages
          .map((message) => message.senderInboxId)
          .filter((value) => typeof value === "string" && value.length > 0),
      ),
    );
    if (client.inboxId) {
      inboxIds.push(client.inboxId);
    }

    const inboxWalletMap = await resolveInboxWalletMap(client, inboxIds);
    const normalizedWalletAddress = String(walletAddress || "").toLowerCase();
    const normalizedRecipient = recipientWalletAddress.toLowerCase();

    const mappedMessages: ConversationDetails["messages"] = rawMessages
      .map((message) => {
        const senderWalletFromInbox = inboxWalletMap.get(message.senderInboxId) || "";
        const isOutgoing =
          message.senderInboxId === client.inboxId
          || senderWalletFromInbox === normalizedWalletAddress;
        const senderWalletAddress = isOutgoing
          ? normalizedWalletAddress
          : senderWalletFromInbox || normalizedRecipient;
        const recipientWalletAddressResolved = isOutgoing
          ? normalizedRecipient
          : normalizedWalletAddress;
        const content =
          typeof message.content === "string"
            ? message.content
            : message.fallback || "[Unsupported XMTP content]";

        return {
          id: message.id,
          conversationId,
          tokenId,
          senderWalletAddress,
          recipientWalletAddress: recipientWalletAddressResolved,
          content,
          transport: "xmtp",
          xmtpMessageId: message.id,
          createdAt: message.sentAt,
          readBy: isOutgoing ? [normalizedWalletAddress] : [],
        };
      });

    mappedMessages.sort(
      (left, right) =>
        new Date(left.createdAt).getTime() - new Date(right.createdAt).getTime(),
    );
    setConversationMessages(mappedMessages);
  };

  const workflowProgress = useMemo(() => {
    const houseIsMinted = Boolean(house);
    const hasListing = Boolean(listing && listing.listingType !== "none");
    const ownerAddress = house?.ownerAddress?.toLowerCase() ?? "";
    const originalOwner = house?.originalOwner?.toLowerCase() ?? "";
    const ownerChanged =
      ownerAddress.length > 0 &&
      originalOwner.length > 0 &&
      ownerAddress !== originalOwner;
    const hasRental = Boolean(house?.rental?.isActive);
    const soldOrRented = ownerChanged || hasRental;
    const hasPaidBill = bills.some((bill) => bill.isPaid);
    const settled = ownerChanged || (hasRental && hasPaidBill);

    const currentStepIndex = settled
      ? 3
      : soldOrRented
        ? 2
        : hasListing
          ? 1
          : 0;

    const steps: Array<{ label: string; status: WorkflowStepStatus }> = [
      {
        label: "Minted",
        status:
          currentStepIndex === 0
            ? "current"
            : houseIsMinted
              ? "completed"
              : "upcoming",
      },
      {
        label: "Listed",
        status:
          currentStepIndex === 1
            ? "current"
            : currentStepIndex > 1
              ? "completed"
              : "upcoming",
      },
      {
        label: "Sold / Rented",
        status:
          currentStepIndex === 2
            ? "current"
            : currentStepIndex > 2
              ? "completed"
              : "upcoming",
      },
      {
        label: "Settled",
        status: currentStepIndex === 3 ? "completed" : "upcoming",
      },
    ];

    let nextAction = "Publish listing terms to make this property actionable.";
    if (currentStepIndex === 1) {
      nextAction =
        "Run buy or rent execution from Marketplace or this property panel.";
    } else if (currentStepIndex === 2) {
      nextAction = ownerChanged
        ? "Ownership transferred onchain. Keep documents and billing synced."
        : "Complete settlement by paying rental bills and finalizing operations.";
    } else if (currentStepIndex === 3) {
      nextAction =
        "Lifecycle complete. You can relist anytime for a new cycle.";
    }

    const quickActions: WorkflowQuickAction[] = [];

    if (currentStepIndex === 0) {
      if (isOwner) {
        quickActions.push({
          label: "Go to List",
          to: `/houses/${tokenId}/list`,
          primary: true,
        });
      }
      quickActions.push({ label: "Go to Marketplace", to: "/marketplace" });
      quickActions.push({ label: "Go to Bills", tab: "bills" });
    } else if (currentStepIndex === 1) {
      quickActions.push({
        label: "Go to Marketplace",
        to: "/marketplace",
        primary: true,
      });
      if (isOwner) {
        quickActions.push({
          label: "Update listing",
          to: `/houses/${tokenId}/list`,
        });
      }
      quickActions.push({ label: "Go to Claim Key", to: "/claim" });
    } else if (currentStepIndex === 2) {
      quickActions.push({ label: "Go to Bills", tab: "bills", primary: true });
      quickActions.push({ label: "Go to Claim Key", to: "/claim" });
      if (isOwner) {
        quickActions.push({
          label: "Go to List",
          to: `/houses/${tokenId}/list`,
        });
      }
    } else {
      if (isOwner) {
        quickActions.push({
          label: "Relist property",
          to: `/houses/${tokenId}/list`,
          primary: true,
        });
      }
      quickActions.push({ label: "Go to Bills", tab: "bills" });
      quickActions.push({ label: "Go to Marketplace", to: "/marketplace" });
    }

    return {
      steps,
      nextAction,
      quickActions,
    };
  }, [bills, house, isOwner, listing, tokenId]);

  useEffect(() => {
    loadHouseData();
  }, [tokenId]);

  const loadHouseData = async () => {
    try {
      setLoading(true);
      setError(null);

      const houseResponse = await apiClient.getHouse(tokenId);
      if (houseResponse.success && houseResponse.data) {
        setHouse(houseResponse.data);
        setListing(houseResponse.data.listing || null);
      }

      const billsResponse = await apiClient.getBills(tokenId);
      if (billsResponse.success && billsResponse.data) {
        setBills(billsResponse.data);
      }
    } catch (err: any) {
      setError(err.message || "Failed to load house data");
    } finally {
      setLoading(false);
    }
  };

  const loadConversations = async () => {
    if (!walletAddress) {
      setConversations([]);
      setActiveConversationId("");
      setConversationMessages([]);
      return;
    }

    const now = new Date();
    const syntheticConversations: ConversationSummary[] = messageRecipients.map(
      (recipient) => ({
        id: `wallet:${recipient.address}`,
        tokenId,
        houseId: house?.houseId || `Asset #${tokenId}`,
        participants: [],
        counterpartWalletAddress: recipient.address,
        counterpartRole: (recipient.role.toLowerCase() as ConversationSummary["counterpartRole"]) || "",
        unreadCount: 0,
        createdAt: now,
        updatedAt: now,
        lastMessageAt: undefined,
        lastMessagePreview: "",
        lastMessage: null,
      }),
    );

    try {
      const response = await apiClient.getConversations(tokenId);
      if (!response.success || !response.data || response.data.length === 0) {
        setConversations(syntheticConversations);
        if (syntheticConversations.length > 0 && !activeConversationId) {
          setActiveConversationId(syntheticConversations[0].id);
        }
        return;
      }
      const backendConversations = response.data;
      const mergedByCounterpart = new Map<string, ConversationSummary>();
      for (const conversation of backendConversations) {
        mergedByCounterpart.set(
          String(conversation.counterpartWalletAddress || "").toLowerCase(),
          conversation,
        );
      }
      for (const syntheticConversation of syntheticConversations) {
        const key = syntheticConversation.counterpartWalletAddress.toLowerCase();
        if (!mergedByCounterpart.has(key)) {
          mergedByCounterpart.set(key, syntheticConversation);
        }
      }
      const nextConversations = Array.from(mergedByCounterpart.values());
      setConversations(nextConversations);

      if (nextConversations.length === 0) {
        setActiveConversationId("");
        setConversationMessages([]);
        return;
      }

      const hasCurrentConversation = nextConversations.some(
        (conversation) => conversation.id === activeConversationId,
      );
      if (!hasCurrentConversation) {
        setActiveConversationId(nextConversations[0].id);
      }
    } catch {
      setConversations(syntheticConversations);
      if (syntheticConversations.length > 0 && !activeConversationId) {
        setActiveConversationId(syntheticConversations[0].id);
      }
    }
  };

  const resolveConversationRecipient = async (
    conversationId: string,
  ): Promise<string> => {
    if (conversationId.startsWith("wallet:")) {
      return conversationId.replace("wallet:", "").toLowerCase();
    }

    try {
      const response = await apiClient.getConversation(conversationId);
      if (response.success && response.data) {
        const counterpartAddress = String(
          response.data.conversation?.counterpartWalletAddress || "",
        ).trim();
        if (counterpartAddress) {
          setConversations((previous) =>
            previous.map((conversation) =>
              conversation.id === conversationId
                ? { ...conversation, unreadCount: 0 }
                : conversation,
            ),
          );
          return counterpartAddress.toLowerCase();
        }
      }
    } catch {
      // fallback to local conversation cache
    }

    const localConversation = conversations.find(
      (conversation) => conversation.id === conversationId,
    );
    return String(localConversation?.counterpartWalletAddress || "")
      .trim()
      .toLowerCase();
  };

  const ensureRoleGatedRecipient = (recipientWalletAddress: string): boolean => {
    const normalizedRecipient = recipientWalletAddress.toLowerCase();
    return messageRecipients.some(
      (recipient) => recipient.address.toLowerCase() === normalizedRecipient,
    );
  };

  const upsertConversationFromRecipient = (
    recipientWalletAddress: string,
    lastMessagePreview: string,
  ) => {
    const normalizedRecipient = recipientWalletAddress.toLowerCase();
    const existingConversation = conversations.find(
      (conversation) =>
        conversation.counterpartWalletAddress.toLowerCase() === normalizedRecipient,
    );
    if (existingConversation) {
      return existingConversation.id;
    }
    const recipientRole =
      messageRecipients.find(
        (recipient) => recipient.address.toLowerCase() === normalizedRecipient,
      )?.role || "";
    const syntheticConversation: ConversationSummary = {
      id: `wallet:${normalizedRecipient}`,
      tokenId,
      houseId: house?.houseId || `Asset #${tokenId}`,
      participants: [],
      counterpartWalletAddress: normalizedRecipient,
      counterpartRole:
        (recipientRole.toLowerCase() as ConversationSummary["counterpartRole"])
        || "",
      unreadCount: 0,
      createdAt: new Date(),
      updatedAt: new Date(),
      lastMessageAt: new Date(),
      lastMessagePreview,
      lastMessage: null,
    };
    setConversations((previous) => [syntheticConversation, ...previous]);
    return syntheticConversation.id;
  };

  const loadConversationDetails = async (conversationId: string) => {
    if (!conversationId) {
      setConversationMessages([]);
      return;
    }

    try {
      setIsLoadingMessages(true);
      const recipientWalletAddress = await resolveConversationRecipient(conversationId);
      if (!recipientWalletAddress) {
        setConversationMessages([]);
        return;
      }

      if (!ensureRoleGatedRecipient(recipientWalletAddress)) {
        setConversationMessages([]);
        return;
      }

      setSelectedRecipient(recipientWalletAddress);
      await loadMessagesFromXmtp(recipientWalletAddress, conversationId);
      setXmtpStatus("ready");
    } catch (error) {
      setXmtpStatus("error");
      setConversationMessages([]);
    } finally {
      setIsLoadingMessages(false);
    }
  };

  const handleSendMessage = async () => {
    if (!walletAddress) {
      toast.error("Connect a wallet to send private messages");
      return;
    }
    const recipient = selectedRecipient.trim().toLowerCase();
    if (!recipient) {
      toast.error("Choose a recipient first");
      return;
    }
    if (!ensureRoleGatedRecipient(recipient)) {
      toast.error(
        "Messaging is role-gated. Allowed pairs are seller↔buyer and landlord↔renter for active listings.",
      );
      return;
    }
    if (!draftMessage.trim()) {
      toast.error("Write a message first");
      return;
    }

    try {
      setIsSendingMessage(true);
      const client = await getXmtpClient();
      const dmConversation = await getXmtpDmConversation(client, recipient);
      const sentMessage = (
        await dmConversation.sendText(draftMessage.trim())
      ) as unknown;
      const xmtpMessageId =
        typeof sentMessage === "string"
          ? sentMessage
          : sentMessage &&
              typeof sentMessage === "object" &&
              "id" in sentMessage
            ? String((sentMessage as { id: unknown }).id)
            : "";

      const mirrorResponse = await apiClient.sendPrivateMessage({
        tokenId,
        recipientWalletAddress: recipient,
        message: draftMessage.trim(),
        xmtpMessageId,
      });
      if (!mirrorResponse.success) {
        toast.error(
          mirrorResponse.message
            || "Message sent on XMTP but backend notification sync failed.",
        );
      }

      const messagePreview = draftMessage.trim().slice(0, 120);
      setDraftMessage("");
      await loadConversations();
      const conversationId =
        mirrorResponse.data?.conversation?.id ||
        upsertConversationFromRecipient(recipient, messagePreview);
      setActiveConversationId(conversationId);
      await loadConversationDetails(conversationId);
      toast.success("Private XMTP message sent");
    } catch (error) {
      setXmtpStatus("error");
      const message =
        error instanceof Error
          ? error.message
          : "Unable to send private XMTP message";
      toast.error(message);
    } finally {
      setIsSendingMessage(false);
    }
  };

  useEffect(() => {
    void loadConversations();
  }, [tokenId, walletAddress, house?.houseId, messageRecipients.length]);

  useEffect(() => {
    if (!activeConversationId) {
      return;
    }
    void loadConversationDetails(activeConversationId);
  }, [activeConversationId]);

  useEffect(() => {
    if (activeTab !== "messages" || !activeConversationId) {
      return;
    }
    const intervalId = window.setInterval(() => {
      void loadConversationDetails(activeConversationId);
    }, 8000);
    return () => {
      window.clearInterval(intervalId);
    };
  }, [activeConversationId, activeTab]);

  useEffect(() => {
    xmtpReachabilityCacheRef.current.clear();
  }, [walletAddress]);

  useEffect(() => {
    if (selectedRecipient) {
      return;
    }
    const activeConversation = conversations.find(
      (conversation) => conversation.id === activeConversationId,
    );
    const activeCounterpart = String(
      activeConversation?.counterpartWalletAddress || "",
    )
      .trim()
      .toLowerCase();
    if (/^0x[a-fA-F0-9]{40}$/.test(activeCounterpart)) {
      setSelectedRecipient(activeCounterpart);
      return;
    }
    const latestCounterpart = String(
      conversations[0]?.counterpartWalletAddress || "",
    )
      .trim()
      .toLowerCase();
    if (/^0x[a-fA-F0-9]{40}$/.test(latestCounterpart)) {
      setSelectedRecipient(latestCounterpart);
      return;
    }
    if (messageRecipients.length > 0) {
      setSelectedRecipient(messageRecipients[0].address);
    }
  }, [activeConversationId, conversations, messageRecipients, selectedRecipient]);

  useEffect(() => {
    const normalizedWalletAddress = String(walletAddress || "").toLowerCase();
    if (
      xmtpClientRef.current &&
      xmtpWalletAddressRef.current &&
      xmtpWalletAddressRef.current !== normalizedWalletAddress
    ) {
      xmtpClientRef.current.close();
      xmtpClientRef.current = null;
      xmtpWalletAddressRef.current = "";
      setXmtpStatus("idle");
    }
  }, [walletAddress]);

  useEffect(() => {
    return () => {
      if (xmtpClientRef.current) {
        xmtpClientRef.current.close();
      }
      xmtpClientRef.current = null;
      xmtpWalletAddressRef.current = "";
    };
  }, []);

  const handlePayBill = async (billIndex: number) => {
    try {
      const response = await apiClient.payBill({
        tokenId: tokenId,
        billIndex,
        paymentMethod: "crypto",
      });

      if (response.success) {
        toast.success("Payment request submitted");
        loadHouseData();
      } else {
        toast.error(response.message || "Payment failed");
      }
    } catch (err: any) {
      setError(err.message);
      toast.error(err.message || "Payment failed");
    }
  };

  const handleCancelListing = async () => {
    if (!isOwner) {
      toast.error("Only the owner can cancel the listing");
      return;
    }

    if (wrongChain) {
      notifyWrongNetwork();
      return;
    }

    const contractAddress = import.meta.env.VITE_HOUSE_RWA_ADDRESS;
    if (!contractAddress) {
      toast.error("Missing VITE_HOUSE_RWA_ADDRESS in web env config");
      return;
    }

    try {
      setActionLoading("cancel");
      const ethereumProvider = await getEthereumProvider();
      const provider = new ethers.BrowserProvider(ethereumProvider);
      const signer = await provider.getSigner();
      const contract = new ethers.Contract(
        contractAddress,
        ["function cancelListing(uint256 tokenId)"],
        signer,
      );

      const tx = await contract.cancelListing(BigInt(tokenId));
      await tx.wait();

      toast.success("Listing cancelled onchain");
      await loadHouseData();
    } catch (err: any) {
      setError(err?.message || "Failed to cancel listing");
      toast.error(
        err?.shortMessage || err?.message || "Failed to cancel listing",
      );
    } finally {
      setActionLoading(null);
    }
  };

  const handleBuy = async () => {
    if (!walletAddress) {
      toast.error("Connect a wallet to buy this property");
      return;
    }
    if (!house || !listing || listing.listingType !== "for_sale") {
      toast.error("This property is not listed for sale");
      return;
    }
    if (wrongChain) {
      notifyWrongNetwork();
      return;
    }
    if (
      listing.isPrivateSale &&
      listing.allowedBuyer &&
      listing.allowedBuyer.toLowerCase() !== walletAddress.toLowerCase()
    ) {
      toast.error("Private listing: your wallet is not the allowed buyer");
      return;
    }

    try {
      setActionLoading("buy");
      const response = await apiClient.sellHouse({
        action: "sell",
        sellerAddress: house.ownerAddress,
        buyerAddress: walletAddress,
        tokenId,
        price: listing.price,
        buyerPublicKey: "",
        isPrivateSale: listing.isPrivateSale,
      });

      if (response.success) {
        const keyHash = String(response.data?.keyHash || "").trim();
        const keyHashSaved = saveLatestClaimKeyHash(keyHash);
        toast.success(
          response.txHash
            ? `Purchase submitted: ${response.txHash.slice(0, 10)}...${
              keyHashSaved ? " Key hash saved for claim." : ""
            }`
            : `Purchase submitted via CRE${
              keyHashSaved ? " (key hash saved for claim)." : ""
            }`,
        );
        await loadHouseData();
      } else {
        toast.error(response.message || "Failed to submit purchase");
      }
    } catch (err: any) {
      const message = err?.message || "Failed to submit purchase";
      setError(message);
      toast.error(message);
    } finally {
      setActionLoading(null);
    }
  };

  const fundRentalDepositIfNeeded = async (requiredDepositWei: string) => {
    if (!walletAddress) {
      throw new Error("Connect a wallet before funding rental deposit");
    }

    const contractAddress = import.meta.env.VITE_HOUSE_RWA_ADDRESS;
    if (!contractAddress) {
      throw new Error("Missing VITE_HOUSE_RWA_ADDRESS in web env config");
    }

    const ethereumProvider = await getEthereumProvider();
    const provider = new ethers.BrowserProvider(ethereumProvider);
    const signer = await provider.getSigner();
    const contract = new ethers.Contract(
      contractAddress,
      [
        "function depositForRental(uint256 tokenId) payable",
        "function pendingRentalDeposits(uint256 tokenId,address renter) view returns (uint256)",
      ],
      signer,
    );

    const required = BigInt(requiredDepositWei);
    const pendingRaw = await contract.pendingRentalDeposits(
      BigInt(tokenId),
      walletAddress,
    );
    const pending = BigInt(pendingRaw);

    if (pending >= required) {
      return;
    }

    const shortfall = required - pending;
    const tx = await contract.depositForRental(BigInt(tokenId), {
      value: shortfall,
    });
    await tx.wait();
  };

  const handleRent = async () => {
    if (!walletAddress) {
      toast.error("Connect a wallet to rent this property");
      return;
    }
    if (!listing || listing.listingType !== "for_rent") {
      toast.error("This property is not listed for rent");
      return;
    }
    if (wrongChain) {
      notifyWrongNetwork();
      return;
    }

    const duration = Number.parseInt(rentDurationDays, 10);
    if (!Number.isFinite(duration) || duration <= 0 || duration > 3650) {
      toast.error("Duration must be between 1 and 3650 days");
      return;
    }

    let depositWei = listing.price;
    const trimmedDeposit = rentDepositAmount.trim();
    if (trimmedDeposit) {
      try {
        depositWei =
          mode === "easy"
            ? ethers.parseEther(trimmedDeposit).toString()
            : trimmedDeposit;
      } catch {
        toast.error(
          mode === "easy"
            ? "Deposit must be a valid ETH amount (e.g. 0.25)"
            : "Deposit must be a valid wei amount",
        );
        return;
      }
    }

    const payload: RentRequestPayload = {
      action: "rent",
      tokenId,
      renterAddress: walletAddress,
      durationDays: duration,
      monthlyRent: listing.price,
      depositAmount: depositWei,
      renterPublicKey: "",
    };

    try {
      setActionLoading("rent");

      toast.loading("Verifying renter KYC...", { id: "rent-flow" });
      const kycResponse = await apiClient.ensureKYC(walletAddress);
      if (!kycResponse.success) {
        toast.error(kycResponse.message || "Unable to verify renter KYC", {
          id: "rent-flow",
        });
        return;
      }

      toast.loading("Funding rental deposit...", { id: "rent-flow" });
      await fundRentalDepositIfNeeded(depositWei);

      toast.loading("Submitting rental agreement via CRE...", {
        id: "rent-flow",
      });
      const response = await apiClient.rentHouse(payload);

      if (response.success) {
        const accessKeyHash = String(response.data?.accessKeyHash || "").trim();
        const keyHashSaved = saveLatestClaimKeyHash(accessKeyHash);
        toast.success(
          response.txHash
            ? `Rental submitted: ${response.txHash.slice(0, 10)}...${
              keyHashSaved ? " Key hash saved for claim." : ""
            }`
            : `Rental submitted via CRE${
              keyHashSaved ? " (key hash saved for claim)." : ""
            }`,
          { id: "rent-flow" },
        );
        await loadHouseData();
      } else {
        toast.error(response.message || "Failed to submit rental", {
          id: "rent-flow",
        });
      }
    } catch (err: any) {
      const message = err?.message || "Failed to submit rental";
      setError(message);
      toast.error(message, { id: "rent-flow" });
    } finally {
      setActionLoading(null);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[70vh]">
        <div className="relative">
          <div className="w-16 h-16 border-2 border-[#00f3ff] border-t-transparent rounded-full animate-spin"></div>
          <div
            className="absolute inset-0 w-16 h-16 border-2 border-[#b026ff] border-b-transparent rounded-full animate-spin"
            style={{ animationDirection: "reverse", animationDuration: "1.5s" }}
          ></div>
        </div>
      </div>
    );
  }

  if (!house) {
    return (
      <div className="page-shell page-shell-tight">
        <div className="cyber-card border-[#ff3366] bg-[rgba(255,51,102,0.08)] p-6">
          <p className="text-[#ff3366] font-mono">House not found</p>
          <div className="mt-4">
            <Link to="/dashboard" className="cyber-btn cyber-btn-primary">
              Back
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page-shell page-shell-tight workspace-surface">
      {/* Header */}
      <div className="cyber-card hero-panel overflow-hidden mb-8 mx-auto max-w-6xl">
        <div className="p-6 md:p-7 border-b border-[rgba(0,243,255,0.2)]">
          <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-5">
            <div>
              <h1 className="form-title">
                {house.metadata?.address ? (
                  house.metadata.address
                ) : (
                  <>
                    Token{" "}
                    <span className="number-pill number-pill-sm number-pill-mono">
                      #{tokenId}
                    </span>
                  </>
                )}
              </h1>
              <p className="form-subtitle font-mono">
                {house.metadata.city}, {house.metadata.state}{" "}
                <span className="number-pill number-pill-sm number-pill-mono">
                  {house.metadata.zipCode}
                </span>
              </p>
              <div className="mt-3 flex flex-wrap gap-2">
                <span className="meta-chip">
                  Token{" "}
                  <span className="number-pill number-pill-sm number-pill-mono">
                    #{tokenId}
                  </span>
                </span>
                <span className="meta-chip">
                  Owner {house.ownerAddress.slice(0, 6)}...
                  {house.ownerAddress.slice(-4)}
                </span>
              </div>
            </div>
            {isOwner && (
              <div className="flex flex-wrap gap-2">
                <Link
                  to={`/houses/${tokenId}/list`}
                  className="cyber-btn cyber-btn-primary"
                >
                  {mode === "degen" ? "List" : "Sell / Rent"}
                </Link>
                <Link
                  to={`/houses/${tokenId}/bills/create`}
                  className="cyber-btn"
                >
                  {mode === "degen" ? "Create Bill" : "Add Utility Bill"}
                </Link>
              </div>
            )}
          </div>

          <div className="mt-5 rounded-2xl border border-[#3b82f640] bg-[#0d1b34b8] p-5 sm:p-6">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-sm font-semibold text-[#dce8ff]">
                CRE workflow progress
              </h2>
              <span className="text-xs text-[#8fb4ff]">
                Minted → Listed → Sold/Rented → Settled
              </span>
            </div>
            <div className="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
              {workflowProgress.steps.map((step, index) => {
                const style = WORKFLOW_STEP_STYLES[step.status];
                const isLastStep = index === workflowProgress.steps.length - 1;

                return (
                  <div
                    key={step.label}
                    className={`relative rounded-xl border p-4 text-left transition ${style.card}`}
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span
                        className={`inline-flex h-8 w-8 items-center justify-center rounded-full border text-xs font-semibold ${style.index}`}
                      >
                        {index + 1}
                      </span>
                      <span
                        className={`inline-flex rounded-full border px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.08em] ${style.badge}`}
                      >
                        {style.label}
                      </span>
                    </div>
                    <p className="mt-3 text-sm font-semibold text-[#e6efff]">
                      {step.label}
                    </p>
                    {!isLastStep && (
                      <span className="pointer-events-none absolute -right-2 top-1/2 hidden h-px w-4 -translate-y-1/2 bg-slate-600/70 xl:block" />
                    )}
                  </div>
                );
              })}
            </div>
            <p className="mt-4 rounded-lg border border-slate-600/45 bg-slate-900/45 px-3 py-2 text-sm text-[#c2d7ff]">
              {workflowProgress.nextAction}
            </p>
            <div className="mt-4 flex flex-wrap gap-2">
              {workflowProgress.quickActions.map((action) => {
                const buttonClass = action.primary
                  ? "cyber-btn cyber-btn-primary text-sm"
                  : "cyber-btn text-sm";

                if (action.to) {
                  return (
                    <Link
                      key={`${action.label}-${action.to}`}
                      to={action.to}
                      className={buttonClass}
                    >
                      {action.label}
                    </Link>
                  );
                }

                if (action.tab) {
                  const tab = action.tab;
                  return (
                    <button
                      key={`${action.label}-${tab}`}
                      type="button"
                      className={buttonClass}
                      onClick={() => setActiveTab(tab)}
                    >
                      {action.label}
                    </button>
                  );
                }

                return null;
              })}
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="border-t border-[rgba(0,243,255,0.15)]">
          <nav className="flex">
            {(["details", "bills", "rental", "messages"] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`detail-tab-btn ${
                  activeTab === tab
                    ? "detail-tab-btn-active"
                    : "detail-tab-btn-idle"
                }`}
              >
                {mode === "easy" && tab === "details"
                  ? "home"
                  : mode === "easy" && tab === "bills"
                    ? "payments"
                    : mode === "easy" && tab === "rental"
                      ? "lease"
                      : mode === "easy" && tab === "messages"
                        ? "chat"
                      : tab}
              </button>
            ))}
          </nav>
        </div>
      </div>

      {/* Content */}
      {error && (
        <div className="cyber-card border-[#ff3366] bg-[rgba(255,51,102,0.08)] text-[#ff3366] px-4 py-3 rounded mb-6 font-mono">
          {error}
        </div>
      )}

      {activeTab === "details" && (
        <div className="mx-auto grid max-w-6xl grid-cols-1 gap-6 md:grid-cols-2">
          {/* Property Info */}
          <div className="cyber-card p-6">
            <h2 className="section-title mb-4">Property Details</h2>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Property Type
                </span>
                <span className="font-medium text-white">
                  {house.metadata.propertyType.replace("_", " ")}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Bedrooms
                </span>
                <span className="number-pill number-pill-sm">
                  {house.metadata.bedrooms}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Bathrooms
                </span>
                <span className="number-pill number-pill-sm">
                  {house.metadata.bathrooms}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Square Feet
                </span>
                <span className="number-pill number-pill-sm number-pill-mono">
                  {house.metadata.squareFeet.toLocaleString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Year Built
                </span>
                <span className="number-pill number-pill-sm">
                  {house.metadata.yearBuilt}
                </span>
              </div>
            </div>
          </div>

          {/* Listing Info */}
          <div className="cyber-card p-6">
            <h2 className="section-title mb-4">
              {mode === "degen" ? "Listing Status" : "Sale / Rent Status"}
            </h2>
            {listing && listing.listingType !== "none" ? (
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-[var(--text-secondary)] font-mono">
                    Status
                  </span>
                  <span
                    className={`font-medium ${
                      listing.listingType === "for_sale"
                        ? "text-[#00ff88]"
                        : "text-[#00f3ff]"
                    }`}
                  >
                    {listing.listingType === "for_sale"
                      ? "For Sale"
                      : "For Rent"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-[var(--text-secondary)] font-mono">
                    Price
                  </span>
                  <span className="number-pill number-pill-sm number-pill-mono">
                    {listing.priceFormatted}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-[var(--text-secondary)] font-mono">
                    Private Sale
                  </span>
                  <span className="font-medium text-white">
                    {listing.isPrivateSale ? "Yes" : "No"}
                  </span>
                </div>
                {listing.isPrivateSale && listing.allowedBuyer && (
                  <div className="flex justify-between">
                    <span className="text-[var(--text-secondary)] font-mono">
                      Allowed Buyer
                    </span>
                    <span className="font-medium text-white font-mono">
                      {listing.allowedBuyer.slice(0, 6)}...
                      {listing.allowedBuyer.slice(-4)}
                    </span>
                  </div>
                )}
                {isOwner && (
                  <button
                    onClick={handleCancelListing}
                    disabled={actionLoading === "cancel" || wrongChain}
                    className="cyber-btn w-full mt-4 border-[#ff3366] text-[#ff3366] hover:bg-[rgba(255,51,102,0.08)] disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {actionLoading === "cancel"
                      ? "Cancelling..."
                      : "Cancel Listing"}
                  </button>
                )}

                {!isOwner && listing.listingType === "for_sale" && (
                  <div className="mt-4 pt-4 border-t border-[rgba(0,243,255,0.15)] space-y-3">
                    <p className="text-xs text-[var(--text-secondary)] font-mono">
                      {mode === "degen"
                        ? "Purchase submits a CRE `sell` action (buyer KYC + `completeSale` write)."
                        : "Buying runs a private compliance check and finalizes ownership onchain."}
                    </p>
                    <button
                      onClick={handleBuy}
                      disabled={
                        actionLoading === "buy" ||
                        wrongChain ||
                        (listing.isPrivateSale &&
                          !!listing.allowedBuyer &&
                          walletAddress?.toLowerCase() !==
                            listing.allowedBuyer.toLowerCase())
                      }
                      className="cyber-btn cyber-btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {actionLoading === "buy"
                        ? mode === "degen"
                          ? "Submitting Buy..."
                          : "Submitting Purchase..."
                        : mode === "degen"
                          ? "Buy via CRE"
                          : "Buy This Property"}
                    </button>
                    {listing.isPrivateSale &&
                      listing.allowedBuyer &&
                      walletAddress?.toLowerCase() !==
                        listing.allowedBuyer.toLowerCase() && (
                        <p className="text-xs text-[#ff3366] font-mono">
                          This private listing is locked to another wallet.
                        </p>
                      )}
                  </div>
                )}

                {!isOwner && listing.listingType === "for_rent" && (
                  <div className="mt-4 pt-4 border-t border-[rgba(0,243,255,0.15)] space-y-3">
                    <p className="text-xs text-[var(--text-secondary)] font-mono">
                      {mode === "degen"
                        ? "Rent funds deposit onchain, then submits CRE `rent` (`startRental`). Deposit defaults to monthly rent if omitted."
                        : "We first lock your deposit in a secure onchain escrow, then create the rental agreement."}
                    </p>
                    <div>
                      <label className="block text-xs text-[var(--text-secondary)] font-mono mb-1">
                        {mode === "degen"
                          ? "Duration (days)"
                          : "Lease Length (days)"}
                      </label>
                      <input
                        className="cyber-input font-mono"
                        type="number"
                        min="1"
                        max="3650"
                        value={rentDurationDays}
                        onChange={(e) => setRentDurationDays(e.target.value)}
                      />
                    </div>
                    <div>
                      <label className="block text-xs text-[var(--text-secondary)] font-mono mb-1">
                        {mode === "degen"
                          ? "Deposit (wei, optional)"
                          : "Deposit (ETH, optional)"}
                      </label>
                      <input
                        className="cyber-input font-mono"
                        type="text"
                        value={rentDepositAmount}
                        onChange={(e) => setRentDepositAmount(e.target.value)}
                        placeholder={
                          mode === "degen"
                            ? listing.price
                            : ethers.formatEther(listing.price)
                        }
                      />
                      <p className="text-xs text-[var(--text-secondary)] mt-1">
                        {mode === "degen" ? (
                          <>
                            Listing rent:{" "}
                            <span className="number-pill number-pill-sm number-pill-mono">
                              {listing.price} wei
                            </span>
                          </>
                        ) : (
                          <>
                            Listing rent:{" "}
                            <span className="number-pill number-pill-sm number-pill-mono">
                              {formatWeiAsEth(listing.price)}
                            </span>
                          </>
                        )}
                      </p>
                    </div>
                    <button
                      onClick={handleRent}
                      disabled={actionLoading === "rent" || wrongChain}
                      className="cyber-btn cyber-btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {actionLoading === "rent"
                        ? mode === "degen"
                          ? "Submitting Rent..."
                          : "Starting Rental..."
                        : mode === "degen"
                          ? "Rent via CRE"
                          : "Start Rental"}
                    </button>
                    <div className="rounded-lg border border-[rgba(0,243,255,0.2)] bg-[rgba(0,243,255,0.07)] p-3 text-xs text-[var(--text-secondary)]">
                      {mode === "degen"
                        ? "Flow: depositForRental -> startRental (CRE) -> claim access key -> pay bills from Payments tab."
                        : "What happens next: 1) deposit is funded, 2) rental is created privately, 3) access key is issued, 4) bills can be paid from the Payments tab."}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <p className="text-[var(--text-secondary)]">
                {mode === "degen"
                  ? "Not currently listed"
                  : "This property is not listed for sale or rent yet"}
              </p>
            )}
          </div>

          {/* Documents */}
          <div className="cyber-card p-6 md:col-span-2">
            <h2 className="section-title mb-4">Documents</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-start">
              <div className="md:col-span-2">
                <p className="text-sm text-[var(--text-secondary)] font-mono">
                  Document Hash
                </p>
                <p className="font-mono text-sm text-white break-all">
                  {house.documentHash}
                </p>
              </div>
              <div>
                <p className="text-sm text-[var(--text-secondary)] font-mono">
                  Storage
                </p>
                <p className="font-medium text-white">{house.storageType}</p>
              </div>
              {isOwner && (
                <div className="md:justify-self-end">
                  <Link
                    to={`/houses/${tokenId}/documents`}
                    className="cyber-btn cyber-btn-primary"
                  >
                    View Documents
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {activeTab === "bills" && (
        <div className="cyber-card mx-auto max-w-6xl overflow-hidden">
          <div className="p-6 md:p-7 border-b border-[rgba(0,243,255,0.2)] flex items-center justify-between gap-4">
            <h2 className="section-title">
              {mode === "degen" ? "Bills" : "Bills & Payments"}
            </h2>
            {isOwner && (
              <Link
                to={`/houses/${tokenId}/bills/create`}
                className="cyber-btn cyber-btn-primary text-sm"
              >
                {mode === "degen" ? "+ Create Bill" : "+ Add Bill"}
              </Link>
            )}
          </div>
          {bills.length === 0 ? (
            <div className="p-10 text-center text-[var(--text-secondary)]">
              No bills yet
            </div>
          ) : (
            <div className="divide-y divide-[rgba(0,243,255,0.12)]">
              {bills.map((bill, index) => (
                <div
                  key={index}
                  className="p-6 flex items-center justify-between gap-6"
                >
                  <div>
                    <p className="font-medium text-white">{bill.billType}</p>
                    <p className="text-sm text-[var(--text-secondary)] font-mono">
                      Amount:{" "}
                      <span className="number-pill number-pill-sm number-pill-mono">
                        {bill.amountFormatted}
                      </span>
                      {" • "}Due:{" "}
                      <span className="number-pill number-pill-sm number-pill-mono">
                        {new Date(bill.dueDate).toLocaleDateString()}
                      </span>
                    </p>
                    <p
                      className={`text-sm mt-1 ${
                        bill.isPaid
                          ? "text-[#00ff88]"
                          : new Date(bill.dueDate) < new Date()
                            ? "text-[#ff3366]"
                            : "text-[#ffaa00]"
                      }`}
                    >
                      Status:{" "}
                      {bill.isPaid
                        ? "paid"
                        : new Date(bill.dueDate) < new Date()
                          ? "overdue"
                          : "pending"}
                    </p>
                  </div>
                  {!bill.isPaid && isOwner && (
                    <div className="flex gap-2">
                      <Link
                        to={`/houses/${tokenId}/pay`}
                        className="cyber-btn text-sm"
                      >
                        {mode === "degen" ? "Pay" : "Pay Bill"}
                      </Link>
                      <button
                        onClick={() => handlePayBill(index)}
                        className="cyber-btn cyber-btn-primary text-sm"
                      >
                        {mode === "degen" ? "Pay (Crypto)" : "Pay with Crypto"}
                      </button>
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {activeTab === "rental" && (
        <div className="cyber-card mx-auto max-w-6xl p-6">
          <h2 className="section-title mb-4">
            {mode === "degen" ? "Rental Status" : "Rental Details"}
          </h2>
          {house.rental && house.rental.isActive ? (
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Renter
                </span>
                <span className="font-medium text-white font-mono">
                  {house.rental.renterAddress}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Start
                </span>
                <span className="font-medium text-white">
                  {new Date(house.rental.startTime).toLocaleDateString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  End
                </span>
                <span className="font-medium text-white">
                  {new Date(house.rental.endTime).toLocaleDateString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-[var(--text-secondary)] font-mono">
                  Deposit
                </span>
                <span className="number-pill number-pill-sm number-pill-mono">
                  {house.rental.depositAmount}
                </span>
              </div>

              <div className="pt-4">
                <p className="text-sm text-[var(--text-secondary)]">
                  {mode === "degen"
                    ? "Rental access keys are delivered via the mediator. If you are the renter, use "
                    : "Your entry key is delivered privately after rental setup. Use "}
                  <Link className="text-[#00f3ff] hover:underline" to="/claim">
                    {mode === "degen" ? "Claim Key" : "Get Access Key"}
                  </Link>
                  {mode === "degen"
                    ? " with the provided keyHash."
                    : " using the key hash you receive."}
                </p>
              </div>
            </div>
          ) : (
            <p className="text-[var(--text-secondary)]">
              {mode === "degen"
                ? "Not currently rented"
                : "No active rental right now"}
            </p>
          )}
        </div>
      )}

      {activeTab === "messages" && (
        <div className="cyber-card mx-auto max-w-6xl overflow-hidden">
          <div className="border-b border-[rgba(0,243,255,0.2)] px-6 py-5 md:px-7">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h2 className="section-title">Private XMTP Messaging</h2>
              <span
                className={`number-pill number-pill-sm ${
                  xmtpStatus === "ready"
                    ? "border border-emerald-300/45 bg-emerald-400/15 text-emerald-200"
                    : xmtpStatus === "connecting"
                      ? "border border-amber-300/45 bg-amber-400/15 text-amber-200"
                      : xmtpStatus === "error"
                        ? "border border-rose-300/45 bg-rose-400/15 text-rose-200"
                        : "border border-slate-500/45 bg-slate-700/20 text-slate-300"
                }`}
              >
                XMTP: {xmtpStatus}
              </span>
            </div>
            <p className="mt-2 text-sm text-[var(--text-secondary)]">
              Messages are sent on the XMTP network and mirrored to backend
              notifications. Role-gating allows only seller↔buyer and
              landlord↔renter pairs for this property.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-0 lg:grid-cols-[300px,1fr]">
            <aside className="border-b border-[rgba(0,243,255,0.12)] lg:border-b-0 lg:border-r lg:border-[rgba(0,243,255,0.12)]">
              <div className="space-y-2 p-4">
                {conversations.length === 0 ? (
                  <p className="rounded-lg border border-dashed border-[rgba(0,243,255,0.2)] px-3 py-4 text-xs text-[var(--text-secondary)]">
                    No conversation yet. Send the first private message.
                  </p>
                ) : (
                  conversations.map((conversation) => (
                    <button
                      key={conversation.id}
                      type="button"
                      onClick={() => setActiveConversationId(conversation.id)}
                      className={`w-full rounded-lg border px-3 py-3 text-left transition ${
                        activeConversationId === conversation.id
                          ? "border-[#00f3ff66] bg-[rgba(0,243,255,0.12)]"
                          : "border-[rgba(0,243,255,0.2)] hover:bg-[rgba(0,243,255,0.08)]"
                      }`}
                    >
                      <p className="text-xs font-semibold uppercase tracking-[0.08em] text-[#8fb4ff]">
                        {conversation.counterpartRole || "counterparty"}
                      </p>
                      <p className="mt-1 text-sm font-medium text-white font-mono">
                        {conversation.counterpartWalletAddress
                          ? `${conversation.counterpartWalletAddress.slice(0, 6)}...${conversation.counterpartWalletAddress.slice(-4)}`
                          : "Unknown"}
                      </p>
                      <p className="mt-2 text-xs text-[var(--text-secondary)] line-clamp-2">
                        {conversation.lastMessagePreview || "No messages yet"}
                      </p>
                      {conversation.unreadCount > 0 && (
                        <span className="mt-2 inline-flex number-pill number-pill-sm number-pill-mono">
                          {conversation.unreadCount} unread
                        </span>
                      )}
                    </button>
                  ))
                )}
              </div>
            </aside>

            <section className="flex min-h-[420px] flex-col">
              <div className="flex-1 space-y-3 overflow-y-auto px-5 py-4 md:px-6">
                {isLoadingMessages ? (
                  <p className="text-sm text-[var(--text-secondary)]">
                    Loading messages...
                  </p>
                ) : conversationMessages.length === 0 ? (
                  <p className="text-sm text-[var(--text-secondary)]">
                    No messages in this conversation yet.
                  </p>
                ) : (
                  conversationMessages.map((message) => {
                    const isOutgoing =
                      message.senderWalletAddress.toLowerCase()
                      === String(walletAddress || "").toLowerCase();
                    return (
                      <div
                        key={message.id}
                        className={`max-w-[92%] rounded-xl border px-3 py-2 ${
                          isOutgoing
                            ? "ml-auto border-[#00f3ff66] bg-[rgba(0,243,255,0.12)]"
                            : "border-[rgba(148,163,184,0.4)] bg-[rgba(15,23,42,0.65)]"
                        }`}
                      >
                        <p className="text-sm text-white">{message.content}</p>
                        <p className="mt-1 text-[11px] text-[var(--text-secondary)]">
                          {new Date(message.createdAt).toLocaleString()}
                        </p>
                      </div>
                    );
                  })
                )}
              </div>

              <div className="border-t border-[rgba(0,243,255,0.12)] px-5 py-4 md:px-6">
                <div className="grid gap-2 sm:grid-cols-[220px,1fr]">
                  <select
                    value={selectedRecipient}
                    onChange={(event) => setSelectedRecipient(event.target.value)}
                    className="cyber-input !h-[44px] !min-h-[44px]"
                  >
                    <option value="">Choose recipient</option>
                    {messageRecipients.map((recipient) => (
                      <option key={recipient.address} value={recipient.address}>
                        {recipient.role}: {recipient.address.slice(0, 6)}...
                        {recipient.address.slice(-4)}
                      </option>
                    ))}
                  </select>
                  <input
                    className="cyber-input"
                    type="text"
                    value={draftMessage}
                    onChange={(event) => setDraftMessage(event.target.value)}
                    placeholder="Send a private XMTP message..."
                  />
                </div>
                <div className="mt-3 flex flex-wrap items-center justify-between gap-2">
                  <p className="text-xs text-[var(--text-secondary)]">
                    Messages create bell notifications for recipients.
                  </p>
                  <button
                    type="button"
                    onClick={handleSendMessage}
                    disabled={isSendingMessage || messageRecipients.length === 0}
                    className="cyber-btn cyber-btn-primary disabled:opacity-50"
                  >
                    {isSendingMessage ? "Sending..." : "Send message"}
                  </button>
                </div>
              </div>
            </section>
          </div>
        </div>
      )}
    </div>
  );
};
