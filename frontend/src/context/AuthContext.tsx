import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react';
import { login as apiLogin } from '../api/client';

interface AuthState {
  token: string | null;
  username: string | null;
}

interface AuthContextType extends AuthState {
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

function decodePayload(token: string): { username: string; exp: number } | null {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [auth, setAuth] = useState<AuthState>(() => {
    const token = localStorage.getItem('token');
    if (token) {
      const payload = decodePayload(token);
      if (payload && payload.exp * 1000 > Date.now()) {
        return { token, username: payload.username };
      }
      localStorage.removeItem('token');
    }
    return { token: null, username: null };
  });

  useEffect(() => {
    if (!auth.token) return;
    const payload = decodePayload(auth.token);
    if (payload && payload.exp * 1000 <= Date.now()) {
      localStorage.removeItem('token');
      setAuth({ token: null, username: null });
    }
  }, [auth.token]);

  const login = useCallback(async (username: string, password: string) => {
    const token = await apiLogin(username, password);
    localStorage.setItem('token', token);
    const payload = decodePayload(token);
    setAuth({ token, username: payload?.username || username });
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setAuth({ token: null, username: null });
  }, []);

  return (
    <AuthContext.Provider value={{ ...auth, login, logout, isAuthenticated: !!auth.token }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
