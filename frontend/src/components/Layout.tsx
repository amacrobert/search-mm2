import type { ReactNode } from 'react';
import { useAuth } from '../context/AuthContext';

export function Layout({ children }: { children: ReactNode }) {
  const { isAuthenticated, username, logout } = useAuth();

  return (
    <div className="app-layout">
      <header className="app-header">
        <h1>Search MM2</h1>
        {isAuthenticated && (
          <div className="header-right">
            <span className="username">{username}</span>
            <button onClick={logout} className="btn btn-secondary btn-sm">
              Log out
            </button>
          </div>
        )}
      </header>
      <main className="app-main">{children}</main>
    </div>
  );
}
