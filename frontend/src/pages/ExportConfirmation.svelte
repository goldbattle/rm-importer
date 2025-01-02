<script lang="ts">
    import { Button, Listgroup, Navbar } from "flowbite-svelte";
    import { filesToExport } from "../logic/checkboxes.svelte";
    import { get } from "svelte/store";

    let exportIds: string[] = get(filesToExport);
</script>
<div>
    <Navbar color="blue" class="sticky top-0">
        <h1 class="font-bold">Confirm the export</h1>
    </Navbar>
    {#if exportIds.length > 0}
        <Listgroup items={exportIds} let:item active={false}>
            {item}
        </Listgroup>
    {/if}
    <div class="fixed bottom-7 right-10">
        <Button pill size="xl">Proceed</Button>
    </div>
</div>