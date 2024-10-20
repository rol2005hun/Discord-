<section class="container">
    <div class="servers">
        {#each servers as server}
            <button class="server" bind:this={serverButton} class:selected={server === selectedServer} 
                on:click={() => selectServer(server)} on:keydown={(e) => e.key === 'Enter' && selectServer(server)}>
                <img src={server.image} alt={server.name} />
            </button>
        {/each}
    </div>

    <div class="channels">
        <div class="channel-info">
            <h3>{selectedServer.name}</h3>
        </div>
        {#each selectedServer.categories as category}
            <div class="category">
                <h4>{category.name}</h4>
                {#each category.channels as channel}
                    <button 
                        class="channel" 
                        class:selected={channel === selectedChannel}
                        on:click={() => selectChannel(channel)}
                    >
                        <span>{channel.name}</span>
                    </button>
                {/each}
            </div>
        {/each}
        <div class="current-user">
            <button on:click={logoutUser} style="border: none;background: transparent;display: flex;align-items: center;justify-content: center;">
                <img src="https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png" alt="You" />
            </button>
            <div class="infos">
                <span>Youdasdasdasdadasdaasd</span>
                <span>Online</span>
            </div>
        </div>
    </div>

    <div class="chatwall">
        <div class="navbar">
            <h2>{selectedChannel.name}</h2>
            <button on:click={() => showInfo = !showInfo}>
                <i class="fa-solid fa-users"></i>
            </button>
        </div>
        <div class="messages">
            {#each messages as message}
                <div class="message">
                    <img src={message.profilePic} alt={message.user} />
                    <div>
                        <strong>{message.user}</strong>
                        <p>{message.text}</p>
                    </div>
                </div>
            {/each}
        </div>
        <div class="message-box">
            <input 
                type="text" 
                bind:value={newMessage} 
                placeholder={`Üzenet: #${selectedChannel.name}`} 
                on:keydown={(e) => e.key === 'Enter' && sendMessage()} 
            />
            <button on:click={sendMessage}>Send</button>
        </div>
    </div>

    {#if showInfo}
        <div class="info">
            <h3>Felhasználók</h3>
            <div class="users">
                {#each users as user}
                    <div class="user">
                        <img src={user.profilePic} alt={user.name} />
                        <span>{user.name}</span>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</section>

<script lang="ts">
    import { goto } from '$app/navigation';
    import { logout } from '$lib/scripts/auth';
    import { addToast } from '$lib/scripts/toastStore';
    import './page.scss';

    let servers = [
        {
            name: 'Server 1',
            image: 'https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png',
            categories: [
                {
                    name: 'Text Channels',
                    channels: [
                        { name: 'General' },
                        { name: 'Random' },
                        { name: 'Support' }
                    ]
                },
                {
                    name: 'Voice Channels',
                    channels: [
                        { name: 'General Voice' },
                        { name: 'Music' }
                    ]
                }
            ]
        },
        {
            name: 'Server 2',
            image: 'https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png',
            categories: [
                {
                    name: 'General',
                    channels: [
                        { name: 'General' },
                        { name: 'Announcements' }
                    ]
                },
                {
                    name: 'Feedback',
                    channels: [
                        { name: 'Feedback' },
                        { name: 'Bugs' }
                    ]
                }
            ]
        }
    ];

    let serverButton: HTMLButtonElement;
    let selectedServer = servers[0];
    let selectedChannel = selectedServer.categories[0].channels[0];
    let showInfo = true;
    let messages = [
        { user: 'rol', text: 'Hello!', profilePic: 'https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png' },
        { user: 'palkatapeti12', text: 'tsa!', profilePic: 'https://th.bing.com/th/id/OIP.A2cmBcRH-eizwzZcF2y9RgHaE7?rs=1&pid=ImgDetMain' },
    ];
    let newMessage = '';

    let users = [
        { name: 'rol', profilePic: 'https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png' },
        { name: 'palkatapeti12', profilePic: 'https://th.bing.com/th/id/OIP.A2cmBcRH-eizwzZcF2y9RgHaE7?rs=1&pid=ImgDetMain' },
    ];

    function sendMessage() {
        if (newMessage.trim()) {
            const message = { user: 'You', text: newMessage, profilePic: 'https://www.mintface.xyz/content/images/2021/08/QmTndiF423kjdXsNzsip1QQkBQqDuzDhJnGuJAXtv4XXiZ-1.png' };
            messages.push(message);
            newMessage = '';
        }
    }

    function selectServer(server: any) {
        selectedServer = server;
        selectedChannel = server.categories[0].channels[0];
    }

    function selectChannel(channel: any) {
        selectedChannel = channel;
    }

    async function logoutUser() {
        try {
            const response = await logout();
            if (response) {
                goto('/');
                addToast('Sikeres kijelentkezés', 'success');
            }
        } catch (error) {
            addToast('Sikertelen kijelentkezés', 'error');
        }
    }
</script>