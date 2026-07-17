export type MedicalDecisionChoice =
  | "rest"
  | "reduce_minutes"
  | "ignore_doctor"
  | "force_return";

export interface MedicalDecisionInput {
  injuryId: string;
  playerId: string;
  choiceId: MedicalDecisionChoice;
}

export interface MedicalDecisionRecord {
  choiceId: MedicalDecisionChoice;
  decisionId: string;
}

export interface MedicalDecisionRequestResult {
  ok: boolean;
  decisionId?: string;
}
