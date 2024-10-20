<svelte:head>
	<title>Auth</title>
	<meta name="description" content="Svelte auth" />
</svelte:head>

<section>
	<div class="auth-page">
		{#if isLogin}
			<div class="login" bind:this={loginDiv}>
				<h2>Bejelentkezés</h2>
				<form class="login-form" in:slide={{ duration: 300 }} out:slide={{ duration: 3000 }} on:submit={handleLogin}>
					<div class="username">
						<label for="usernameOrEmail">Felhasználónév / email:</label>
						<input type="text" id="usernameOrEmail" name="usernameOrEmail" bind:value={usernameOrEmail} required />
					</div>
					<div class="password">
						<label for="password">Jelszó:</label>
						<input type="password" id="password" name="password" bind:value={password} required />
					</div>
					<button type="submit">Bejelentkezés</button>
				</form>
				<p>Nincs még fiókod? <button on:click={toggleForm}>Regisztrálj!</button></p>
			</div>
		{:else}
			<div class="register" bind:this={registerDiv}>
				<h2>Regisztráció</h2>
				<form class="register-form" in:slide={{ duration: 300 }} out:slide={{ duration: 300 }} on:submit={handleRegister}>
					<div class="username">
						<label for="username">Felhasználónév:</label>
						<input type="text" id="username" name="username" bind:value={username} required />
					</div>
					<div class="email">
						<label for="email">Email:</label>
						<input type="email" id="email" name="email" bind:value={email} required />
					</div>
					<div class="password">
						<label for="password">Jelszó:</label>
						<input type="password" id="password" name="password" bind:value={password} required />
					</div>
					<button type="submit">Regisztráció</button>
				</form>
				<p>Már van fiókod? <button on:click={toggleForm}>Jelentkezz be!</button></p>
			</div>
		{/if}
	</div>
</section>

<script lang="ts">
	import { slide } from 'svelte/transition';
	import { login, register } from '$lib/scripts/auth';
	import { addToast } from '$lib/scripts/toastStore';
	import './page.scss';
    import { goto } from '$app/navigation';

	let isLogin = true;
	let loginDiv: HTMLDivElement;
	let registerDiv: HTMLDivElement;

	let usernameOrEmail = '';
	let username = '';
	let email = '';
	let password = '';

	function toggleForm() {
		isLogin = !isLogin;

		if (isLogin) {
            registerDiv.style.display = 'none';
			loginDiv.style.display = 'block';
        } else {
            loginDiv.style.display = 'none';
			registerDiv.style.display = 'block';
        }
	}

	async function handleLogin(event: Event) {
		event.preventDefault();

		try {
			const result = await login(usernameOrEmail, password);
			if (result) {
				goto('/server/1/channel/1');
				addToast('Sikeres bejelentkezés', 'success');
			}
		} catch (error) {
			addToast('Sikertelen bejelentkezés', 'error');
		}
	}

	async function handleRegister(event: Event) {
		event.preventDefault();

		try {
			const result = await register(username, email, password);
			addToast('Sikeres regisztráció', 'success');
			email = '';
			password = '';
		} catch (error) {
			addToast('Sikertelen regisztráció', 'error');
		}
	}
</script>
