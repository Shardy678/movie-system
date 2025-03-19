import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';

interface TokenPayload {
  role: string;
}

export function useAuth() {
  const [isAdmin, setIsAdmin] = useState<boolean>(false);
  const navigate = useNavigate();
  const token = localStorage.getItem('token');

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }

    try {
      const decodedToken: TokenPayload = jwtDecode(token);
      setIsAdmin(decodedToken.role === 'admin');
    } catch {
      navigate('/login');
    }
  }, [navigate, token]);

  return { isAdmin, token };
}
