import { redirect } from '@sveltejs/kit';

export function load({ cookies }) {
    const token = cookies.get('token');

    if (token) {
        throw redirect(302, '/server/1/channel/1');
    }

    return {};
}