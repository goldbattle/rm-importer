<script lang="ts">
    import { Alert, Button, Listgroup, Navbar, P, Spinner } from "flowbite-svelte";
    import { Export, GetCheckedFiles} from '../../wailsjs/go/main/App.js';
    import { CheckOutline, ExclamationCircleOutline, FileLinesSolid, InfoCircleSolid } from "flowbite-svelte-icons";
    import { EventsOn } from "../../wailsjs/runtime/runtime.js";
    import { backend } from "../../wailsjs/go/models.js";
    type DocInfo = backend.DocInfo;

    let exportItems: DocInfo[] = $state([]);
    let exportItemState: {[key: string]: string;} = $state({});

    let errorMessage: string = $state("hello");
    let showError: boolean = $state(false);
    let errorHappened: boolean = $state(false);

    GetCheckedFiles()
        .then((result: DocInfo[]) => {
            exportItems = result;
        });
    
    Export();

    const onRetry = () => {
        errorHappened = false;
        Export();
    };

    const finishedAllItems = $derived(() => {
        return Object.keys(exportItemState)
            .filter((id: string) => exportItemState[id] == "finished")
            .length == exportItems.length
    });

    EventsOn("downloading", (id: string) => {
        exportItemState[id] = "downloading";
    });

    EventsOn("finished", (id: string) => {
        exportItemState[id] = "finished";
    });

    EventsOn("error", (id: string, msg: string) => {
        exportItemState[id] = "error";
        showError = true;
        errorMessage = msg;
        errorHappened = true;
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

    <Navbar color="blue" class="sticky top-0">
        <h1 class="font-bold m-auto">{finishedAllItems() ? "Success!" : "Export"}</h1>
    </Navbar>

    <main class="pr-7 pl-7 mt-3 w-full">
    {#if exportItems.length > 0}
        <Listgroup items={exportItems} let:item active={false}>
            <div class="flex flex-row justify-start items-center w-full">
                <FileLinesSolid class="mr-1" size="lg" />
                <P size="xl">{item.Path}</P>
                {#if exportItemState[item.Id] === "downloading"}
                    <Spinner class="ml-auto" />
                {:else if exportItemState[item.Id] === "finished"}
                    <CheckOutline class="ml-auto" color="green"/>
                {:else if exportItemState[item.Id] === "error"}
                    <ExclamationCircleOutline class="ml-auto" color="red"/>
                {/if}
            </div>
        </Listgroup>
    {/if}
    </main>

    <div class="fixed bottom-7 right-10">
        <Button class={!errorHappened ? "invisible": ""} pill size="xl" onclick={onRetry}>Retry</Button>
    </div>
</div>