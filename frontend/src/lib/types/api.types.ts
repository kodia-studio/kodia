export interface ApiResponse<T = any> {
	success: boolean;
	message: string;
	data?: T;
	errors?: Record<string, string[]>;
	meta?: PaginationMeta;
}

export interface PaginationMeta {
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
}
