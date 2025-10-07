import { apiGet, apiPost, apiPatch, apiDel, useApiClient } from "./_client";
import type { Page } from "./_base";

export interface Member {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  role?: "管理员" | "编辑者" | "查看者" | "协作者";
  status?: "active" | "pending" | "disabled";
  joined_at?: string;
  last_active?: string;
  permissions?: string[];
}

export interface MemberInvite {
  email: string;
  role: string;
  message?: string;
}

export function useMemberApi() {
  const { baseURL } = useApiClient();

  const listMembers = (
    page = 1,
    page_size = 20,
    filters?: {
      role?: string;
      status?: string;
      search?: string;
    }
  ) => {
    const query: Record<string, any> = { page, page_size };
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          query[key] = value;
        }
      });
    }
    return apiGet<Page<Member>>("members", query);
  };

  const getMember = (id: string) => apiGet<Member>(`members/${id}`);

  const createMember = (data: Partial<Member>) =>
    apiPost<Member>("members", data);

  const updateMember = (id: string, data: Partial<Member>) =>
    apiPatch<Member>(`members/${id}`, data);

  const deleteMember = (id: string) => apiDel<void>(`members/${id}`);

  const inviteMember = (data: MemberInvite) =>
    apiPost<{ success: boolean; message: string }>("members/invite", data);

  const resendInvite = (id: string) =>
    apiPost<{ success: boolean; message: string }>(
      `members/${id}/resend-invite`
    );

  const activateMember = (id: string) =>
    apiPatch<Member>(`members/${id}`, { status: "active" });

  const deactivateMember = (id: string) =>
    apiPatch<Member>(`members/${id}`, { status: "disabled" });

  const updateMemberRole = (id: string, role: string) =>
    apiPatch<Member>(`members/${id}`, { role });

  const getMemberPermissions = (id: string) =>
    apiGet<{ permissions: string[] }>(`members/${id}/permissions`);

  const updateMemberPermissions = (id: string, permissions: string[]) =>
    apiPatch<Member>(`members/${id}/permissions`, { permissions });

  const searchMembers = (query: string, limit = 10) =>
    apiGet<Member[]>("members/search", { q: query, limit });

  const getMemberActivity = (id: string, days = 30) =>
    apiGet<{
      login_count: number;
      templates_created: number;
      templates_updated: number;
      last_login: string;
    }>(`members/${id}/activity`, { days });

  return {
    baseURL,
    listMembers,
    getMember,
    createMember,
    updateMember,
    deleteMember,
    inviteMember,
    resendInvite,
    activateMember,
    deactivateMember,
    updateMemberRole,
    getMemberPermissions,
    updateMemberPermissions,
    searchMembers,
    getMemberActivity,
  };
}
