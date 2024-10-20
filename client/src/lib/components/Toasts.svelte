<script lang="ts">
    import { subscribe, removeToast } from '$lib/scripts/toastStore';
    import { onDestroy } from 'svelte';

    let toasts: any[] = [];

    const unsubscribe = subscribe(value => {
        toasts = value;
    });

    onDestroy(() => {
        unsubscribe();
    });
</script>

<style lang="scss">
    .toasts {
        position: fixed;
        top: 0.3em;
        right: 0.3em;
        z-index: 1000;
    }

    .toast {
        padding: 0.8em;
        margin: 10px;
        width: 20em;
        border-radius: 5px;
        position: relative;
        transition: opacity 0.3s ease;
        word-break: break-all;

        span {
            display: block;
            width: 90%;
        }
    }

    .success {
        background-color: #dff0d8;
        color: #3c763d;
    }

    .error {
        background-color: #f2dede;
        color: #a94442;
    }

    .warning {
        background-color: #fcf8e3;
        color: #8a6d3b;
    }

    .close {
        position: absolute;
        right: 0.7em;
        top: 0.35em;
        cursor: pointer;
        border: none;
        background-color: transparent;
        font-size: 1.3em;
    }
</style>

<section class="toasts">
    {#each toasts as toast (toast.id)}
        <div class="toast {toast.type}">
            <span>{toast.message}</span>
            <button class="close" on:click={() => removeToast(toast.id)}>x</button>
        </div>
    {/each}
</section>
