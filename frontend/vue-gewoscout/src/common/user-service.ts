import type { User } from '@/types/User';
import axios from 'axios';

export async function getLoggedInUser(): Promise<User | null> {
  try {
    const response = await axios.get(`/.auth/me`);

    if (response.status !== 200) {
      return null;
    }

    if (!response.data.clientPrincipal) {
      return null;
    }

    return {
      identityProvider: response.data.clientPrincipal.identityProvider,
      userId: response.data.clientPrincipal.userId,
      userDetails: response.data.clientPrincipal.userDetails,
      userRoles: response.data.clientPrincipal.userRoles
    };
  } catch (error) {
    console.error(error);
    return null;
  }
}

export async function logoutUser(): Promise<boolean> {
  try {
    const response = await axios.get(`/.auth/logout`);

    if (response.status !== 200) {
      return false;
    }

    if (!response.data.clientPrincipal) {
      return false;
    }

    return true;
  } catch (error) {
    console.error(error);
    return false;
  }
}
