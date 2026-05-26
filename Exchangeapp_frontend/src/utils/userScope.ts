/** 将超管勾选的用户 ID 追加到查询参数（逗号分隔） */
export function appendSelectedUserIds(params: URLSearchParams, userIds: string[] | null | undefined) {
  if (!userIds || userIds.length === 0) return;
  const valid = userIds.map(String).filter((id) => id && id !== 'null');
  if (valid.length === 0) return;
  params.append('user_ids', valid.join(','));
}

export function hasSelectedUsers(userIds: string[] | null | undefined): boolean {
  return !!userIds && userIds.some((id) => id && id !== 'null');
}

/** 仅选中 1 个用户时返回其 ID */
export function singleSelectedUserId(userIds: string[] | null | undefined): string | null {
  if (!userIds || userIds.length !== 1) return null;
  const id = String(userIds[0]);
  return id && id !== 'null' ? id : null;
}
