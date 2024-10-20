export async function login(usernameOrEmail: string, password: string) {
    try {
        const response = await fetch('http://localhost:8080/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ usernameOrEmail, password }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error during login:', error);
        throw error;
    }
}

export async function register(username: string, email: string, password: string) {
    try {
        const response = await fetch('http://localhost:8080/users/createUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, email, password }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Registration failed');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error during registration:', error);
        throw error;
    }
}

export async function logout() {
    try {
        const response = await fetch('http://localhost:8080/users/logout', {
            method: 'POST',
            credentials: 'include'
        });

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error during logout:', error);
        throw error;
    }
}