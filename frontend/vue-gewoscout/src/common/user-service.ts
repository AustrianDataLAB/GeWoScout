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
