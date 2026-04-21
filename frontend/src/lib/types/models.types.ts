export interface User {
  id: string;
  email: string;
  name: string;
  avatar?: string;
  role: 'user' | 'admin' | 'moderator';
  permissions: string[];
  twoFactorEnabled?: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AuthSession {
  user: User;
  accessToken: string;
  refreshToken: string;
}

export interface Activity {
  id: string;
  userId: string;
  action: string;
  status: string;
  description: string;
  ipAddress: string;
  userAgent: string;
  createdAt: string;
}

export interface DashboardStats {
  totalUsers: number;
  activeSessions: number;
  totalStorage: number;
  errorRate: number;
}
