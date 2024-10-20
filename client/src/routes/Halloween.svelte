<style>
    .pumpkin-button {
        display: flex;
        position: fixed;
        bottom: 0.4em;
        right: 0.4em;
        background-color: orange;
        border: none;
        border-radius: 50%;
        width: 1.5em;
        height: 1.5em;
        cursor: pointer;
        font-size: 3em;
        color: white;
        align-items: center;
        justify-content: center;
    }
</style>

<button class="pumpkin-button" on:click={turnHalloween}>
    <p>🎃</p>
</button>

<script lang="ts">
    import './Halloween.scss';
    let halloweenMode = false;

    function turnHalloween() {
        halloweenMode = !halloweenMode;
        document.body.classList.toggle('halloween', halloweenMode);
    }
</script>