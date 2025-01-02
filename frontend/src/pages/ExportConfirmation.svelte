<script lang="ts">
    import { Button, Listgroup, Navbar, P } from "flowbite-svelte";
    import { filesToExport } from "../logic/checkboxes.svelte";
    import { get } from "svelte/store";
    import { push } from "svelte-spa-router";
    import { GetElementsByIds } from '../../wailsjs/go/main/App.js';
    import { FileLinesSolid } from "flowbite-svelte-icons";

    let exportIds: string[] = get(filesToExport);
    let exportItems: DocInfo[] = $state([]);

    $effect(() => {
        GetElementsByIds(exportIds)
            .then((result: string[]) => {
                exportItems = result;
            });
    });

    const onProceed = () => {
        filesToExport.set(exportIds);
        push('/export');
    };
</script>
<div>
    <Navbar color="blue" class="sticky top-0">
        <h1 class="font-bold m-auto">Confirm the export</h1>
    </Navbar>
    {#if exportItems.length > 0}
        <Listgroup items={exportItems} let:item active={false}>
            <div class="flex flex-row justify-start items-center w-full">
                <FileLinesSolid class="mr-1" size="lg" />
                <P size="xl">{item.Name}</P>
                {#if item.LastModified}
                <P size="md" class="ml-auto">Last modified: {item.LastModified}</P>
                {/if}
            </div>
        </Listgroup>
    {/if}
    <div class="fixed bottom-7 right-10">
        <Button pill size="xl" onclick={onProceed}>Proceed</Button>
    </div>
</div>