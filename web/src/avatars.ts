import type { User } from './types';

export interface AvatarAsset {
  id: string;
  label: string;
}

export type AvatarSelection = { avatar_asset_id: string } | { avatar_url: string };

export const avatarAssetURL = (id: string) => `/api/avatar-catalog/${encodeURIComponent(id)}`;

export function avatarSources(user: Pick<User, 'avatar_asset_id' | 'avatar_url'>) {
  if (!user.avatar_asset_id) return { src: user.avatar_url || null };
  return { src: avatarAssetURL(user.avatar_asset_id), fallbackSrc: user.avatar_url || null };
}
