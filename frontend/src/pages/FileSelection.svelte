<script lang="ts">
    import { Heading, Navbar, P, Toolbar, ToolbarButton } from "flowbite-svelte";
    import { ArrowLeftOutline, CheckOutline, ImageOutline, MapPinAltSolid, PaperClipOutline } from "flowbite-svelte-icons";
    import { push } from "svelte-spa-router";

    export let params = {}
</script>

<Navbar color="blue">
    <ToolbarButton color="blue" name="Attach file" onclick={() => push('/')}><ArrowLeftOutline class="w-7 h-7" /></ToolbarButton>
    <h1 class="font-bold">Choose files to export</h1>
    <ToolbarButton color="blue" name="Upload image"><CheckOutline class="w-7 h-7" /></ToolbarButton>
</Navbar>