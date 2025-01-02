<script lang="ts">
    import { Alert, Button, Heading, Listgroup, Navbar, P, Spinner } from "flowbite-svelte";
    import { filesToExport } from "../logic/checkboxes.svelte";
    import { get } from "svelte/store";
    import { GetElementsByIds, ExportPdfs } from '../../wailsjs/go/main/App.js';
    import { CheckOutline, ExclamationCircleOutline, FileLinesSolid, InfoCircleSolid } from "flowbite-svelte-icons";
    import { EventsOn } from "../../wailsjs/runtime";

    let exportIds: string[] = get(filesToExport);
    let exportItems: DocInfo[] = $state([]);
    let exporting: boolean = $state(false);
    let exportState: {[key: string]: string;} = $state({});
    let errorMessage: string = $state("hello");
    let showError: boolean = $state(false);
    let showExportMessage = $state(true);

    $effect(() => {
        GetElementsByIds(exportIds)
            .then((result: DocInfo[]) => {
                exportItems = result;
            });
    });

    const onProceed = () => {
        showError = false;
        showExportMessage = false;
        exporting = true;
        ExportPdfs(exportIds);
    };

    EventsOn("downloading", (id: string) => {
        exportState[id] = "downloading";
    });

    EventsOn("finished", (id: string) => {
        exportState[id] = "finished";
    });

    EventsOn("error", (id: string, msg: string) => {
        exportState[id] = "error";
        showError = true;
        errorMessage = msg;
        exporting = false;
    });
</script>

<div style="height: fit-content;">
    {#key showError}
    {#if showError}
    <Alert color="red" dismissable={true}>
        <span slot="icon">
            <InfoCircleSolid class="w-5 h-5" />
            <span class="sr-only">Info</span>
        </span>
        <p class="font-large">Couldn't export a file!</p>
        <p class="font-medium">{errorMessage}</p>
    </Alert>
    {/if}
    {/key}

    {#if showExportMessage}
        <Navbar color="blue" class="sticky top-0">
            {#if exporting}
            <h1 class="font-bold m-auto">Following items will be exported in .pdf format:</h1>
            {/if}
        </Navbar>
    {/if}
    {#if exportItems.length > 0}
        <Listgroup items={exportItems} let:item active={false}>
            <div class="flex flex-row justify-start items-center w-full">
                <FileLinesSolid class="mr-1" size="lg" />
                <P size="xl">{item.Name}</P>
                {#if exportState[item.Id] === "downloading"}
                    <Spinner class="ml-auto" />
                {:else if exportState[item.Id] === "finished"}
                    <CheckOutline class="ml-auto" color="green"/>
                {:else if exportState[item.Id] === "error"}
                    <ExclamationCircleOutline class="ml-auto" color="red"/>
                {/if}
            </div>
        </Listgroup>
    {/if}
    <div class="fixed bottom-7 right-10">
        <Button disabled={exporting} pill size="xl" onclick={onProceed}>Proceed</Button>
    </div>
</div>