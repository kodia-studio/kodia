export interface UserResponse {
	id: string;
	name: string;
	email: string;
	role: 'admin' | 'user';
	is_active: boolean;
	avatar_url: string | null;
	created_at: string;
	updated_at: string;
}

export interface UpdateUserRequest {
	name?: string;
	avatar_url?: string;
}

export interface ChangePasswordRequest {
	current_password: string;
	new_password: string;
}
