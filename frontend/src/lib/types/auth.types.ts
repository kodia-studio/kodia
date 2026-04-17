import type { UserResponse } from './user.types';

export interface AuthResponse {
	access_token: string;
	refresh_token: string;
	token_type: string;
	user: UserResponse;
}

export interface RegisterRequest {
	name: string;
	email: string;
	password: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RefreshTokenRequest {
	refresh_token: string;
}
