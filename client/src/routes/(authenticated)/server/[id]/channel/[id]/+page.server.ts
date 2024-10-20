import { redirect } from '@sveltejs/kit';

export async function load({ cookies }) {
    const token = cookies.get('token');

    if (!token) {
        throw redirect(302, '/');
    } else {
        try {
            var result = await fetch('http://localhost:8080/servers/getUserServers', {
                method: 'GET',
                credentials: 'include'
            });
            console.log(result);
        } catch (error) {
            console.error('Error during server load:', error);
            throw error;
        }
    }

    return {servers: result};
}