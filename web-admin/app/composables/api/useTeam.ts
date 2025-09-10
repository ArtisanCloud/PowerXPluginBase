import { apiGet, apiPost, apiPatch, apiDel, useApiClient } from "./_client";

export interface Team {
  id: string;
  name: string;
  description?: string;
  type?: "个人" | "团队" | "企业";
  lead_id?: string;
  lead_name?: string;
  member_count?: number;
  notes_count?: number;
  created_at: string;
  updated_at: string;
  settings?: {
    is_public: boolean;
    allow_member_invite: boolean;
    default_note_permission: string;
  };
}

export interface TeamStats {
  totalMembers: number;
  activeMembers: number;
  totalNotes: number;
  publishedNotes: number;
  draftNotes: number;
  activityRate: number;
  newMembersThisMonth: number;
  notesCreatedThisMonth: number;
}

export interface TeamActivity {
  date: string;
  notes_created: number;
  notes_updated: number;
  members_joined: number;
  active_members: number;
}

export function useTeamApi() {
  const { baseURL } = useApiClient();

  const listTeams = (filters?: { type?: string; search?: string }) => {
    const query: Record<string, any> = {};
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          query[key] = value;
        }
      });
    }
    return apiGet<Team[]>("teams", query);
  };

  const getTeam = (id: string) => apiGet<Team>(`teams/${id}`);

  const getCurrentTeam = () => apiGet<Team>("teams/current");

  const createTeam = (data: Partial<Team>) => apiPost<Team>("teams", data);

  const updateTeam = (id: string, data: Partial<Team>) =>
    apiPatch<Team>(`teams/${id}`, data);

  const deleteTeam = (id: string) => apiDel<void>(`teams/${id}`);

  const getTeamStats = (id?: string) => {
    const path = id ? `teams/${id}/stats` : "teams/current/stats";
    return apiGet<TeamStats>(path);
  };

  const getTeamActivity = (id?: string, days = 30) => {
    const path = id ? `teams/${id}/activity` : "teams/current/activity";
    return apiGet<TeamActivity[]>(path, { days });
  };

  const updateTeamSettings = (
    id: string,
    settings: Partial<Team["settings"]>
  ) => apiPatch<Team>(`teams/${id}/settings`, { settings });

  const transferOwnership = (id: string, newOwnerId: string) =>
    apiPost<{ success: boolean; message: string }>(`teams/${id}/transfer`, {
      new_owner_id: newOwnerId,
    });

  const leaveTeam = (id: string) =>
    apiPost<{ success: boolean; message: string }>(`teams/${id}/leave`);

  const dissolveTeam = (id: string) =>
    apiPost<{ success: boolean; message: string }>(`teams/${id}/dissolve`);

  const getTeamInviteLink = (id: string, expiresIn = 7) =>
    apiPost<{ invite_link: string; expires_at: string }>(
      `teams/${id}/invite-link`,
      { expires_in_days: expiresIn }
    );

  const revokeInviteLink = (id: string) =>
    apiDel<{ success: boolean }>(`teams/${id}/invite-link`);

  return {
    baseURL,
    listTeams,
    getTeam,
    getCurrentTeam,
    createTeam,
    updateTeam,
    deleteTeam,
    getTeamStats,
    getTeamActivity,
    updateTeamSettings,
    transferOwnership,
    leaveTeam,
    dissolveTeam,
    getTeamInviteLink,
    revokeInviteLink,
  };
}
