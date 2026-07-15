export interface TradeProposalInput {
  rivalTeamId: string;
  offeredPlayerId: string;
  requestedPosition: string;
  incomingSalary: number;
}

export interface TradeAcceptanceInput {
  proposalId: string;
  acceptedAdditionalAsset: string;
}

export interface TradeRequestResult {
  ok: boolean;
  proposalId?: string;
}
