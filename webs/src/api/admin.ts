import request from "@/utils/request";

export interface AdminPullLogItem {
  ID?: number;
  IP?: string;
  Date?: string;
  Addr?: string;
  Region?: string;
  Client?: string;
  Status?: string;
  Count?: number;
  SubcriptionID?: number;
  UserID?: number;
  Username?: string;
}

export interface AdminUserItem {
  id: number;
  username: string;
  nickname: string;
  role: string;
  subscriptionId: number;
  subscriptionToken: string;
  allowedRegions: string;
  inviteCodeId: number;
  disabled: boolean;
  createdAt: string;
  pullLogs?: AdminPullLogItem[];
}

export interface InviteItem {
  id: number;
  code: string;
  description: string;
  enabled: boolean;
  usedCount: number;
  createdAt: string;
}

export interface AdminConfig {
  jwt_secret?: string;
  api_encryption_key?: string;
  expire_days?: number;
  port?: number;
  default_subscription_id: number;
  invite_required: boolean;
}

export function getAdminUsers() {
  return request({ url: "/api/v1/admin/users", method: "get" });
}

export function updateAdminUser(data: { id: number; role?: string; subscriptionId?: number; allowedRegions?: string; subscriptionToken?: string; disabled?: boolean }) {
  const form = new FormData();
  form.append("id", String(data.id));
  if (data.role !== undefined) form.append("role", data.role);
  if (data.subscriptionId !== undefined) form.append("subscriptionId", String(data.subscriptionId));
  if (data.allowedRegions !== undefined) form.append("allowedRegions", data.allowedRegions);
  if (data.subscriptionToken !== undefined) form.append("subscriptionToken", data.subscriptionToken);
  if (data.disabled !== undefined) form.append("disabled", String(data.disabled));
  return request({
    url: "/api/v1/admin/users/update",
    method: "post",
    data: form,
    headers: { "Content-Type": "multipart/form-data" }
  });
}

export function resetSubscriptionToken(id: number) {
  const form = new FormData();
  form.append("id", String(id));
  return request({
    url: "/api/v1/admin/users/reset-subscription-token",
    method: "post",
    data: form,
    headers: { "Content-Type": "multipart/form-data" }
  });
}

export function deleteAdminUser(id: number) {
  return request({
    url: `/api/v1/admin/users/${id}`,
    method: "delete",
  });
}

export function getInvites() {
  return request({ url: "/api/v1/admin/invites", method: "get" });
}

export function addInvite(data: { code?: string; description?: string }) {
  const form = new FormData();
  if (data.code) form.append("code", data.code);
  if (data.description) form.append("description", data.description);
  return request({
    url: "/api/v1/admin/invites/add",
    method: "post",
    data: form,
    headers: { "Content-Type": "multipart/form-data" }
  });
}

export function updateInvite(data: { id: number; description?: string; enabled?: boolean }) {
  const form = new FormData();
  form.append("id", String(data.id));
  if (data.description !== undefined) form.append("description", data.description);
  if (data.enabled !== undefined) form.append("enabled", String(data.enabled));
  return request({
    url: "/api/v1/admin/invites/update",
    method: "post",
    data: form,
    headers: { "Content-Type": "multipart/form-data" }
  });
}

export function getAdminConfig() {
  return request({ url: "/api/v1/admin/config", method: "get" });
}

export function setAdminConfig(data: { defaultSubscriptionId?: number; inviteRequired?: boolean }) {
  const form = new FormData();
  if (data.defaultSubscriptionId !== undefined) form.append("defaultSubscriptionId", String(data.defaultSubscriptionId));
  if (data.inviteRequired !== undefined) form.append("inviteRequired", String(data.inviteRequired));
  return request({
    url: "/api/v1/admin/config",
    method: "post",
    data: form,
    headers: { "Content-Type": "multipart/form-data" }
  });
}
