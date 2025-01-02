<script lang="ts">
    import { Button, Checkbox, Listgroup, Navbar, P, ToolbarButton } from "flowbite-svelte";
    import { ArrowUpOutline, FileLinesSolid, FolderSolid } from "flowbite-svelte-icons";
    import { GetTabletFolder } from "../../wailsjs/go/main/App";
    import { isChecked, setChecked, isExportButtonDisabled, storeCheckedFiles } from "../logic/checkboxes.svelte";
    import { push } from "svelte-spa-router";
    
    let id = $state("");
    let path: string[] = $state([]);
    let items: DocInfo[] = $state([]);

    const getFolderData = () => {
        GetTabletFolder(id).then((result) => {
            items = result;
        });
    };

    $effect(getFolderData);

    const onBack = () => {
        if (path.length > 0) {
            id = path[path.length - 1];
            path.pop();
        } else {
            id = '';
        }
    };

    const onExportClick = () => {
        storeCheckedFiles();
        push('/export');
    };
</script>

<div style="height: fit-content;">
    <Navbar color="blue" class="sticky top-0">
        <div class={path.length == 0 ? "invisible": ""}>
            <ToolbarButton color="blue" name="Back" onclick={onBack}> <ArrowUpOutline class="w-7 h-7" /></ToolbarButton>
        </div>
        <h1 class="font-bold">Choose files to export</h1>
        <div>
            <ToolbarButton class="invisible"><ArrowUpOutline class="w-7 h-7" /></ToolbarButton>
        </div>
    </Navbar>
    <main class="pl-10 pr-10 pt-3 pb-3">
        {#if items.length > 0}
        <Listgroup {items} let:item active={false}>
            <div class="flex flex-row justify-start items-center">
                <Checkbox checked={isChecked(item.Id)} 
                            on:click={(e) => setChecked(item.Id, e.target.checked)}
                class={item.IsFolder ? "invisible mr-2" : "mr-2"} />

                <div class="flex flex-row justify-start items-center w-full hover:bg-gray-100"
                     onclick={()=>{
                        if (item.IsFolder) {
                            path.push(id);
                            id = item.Id;
                        }
                     }}>
                    {#if item.IsFolder}
                        <FolderSolid class="mr-1" size="lg" />
                    {:else}
                        <FileLinesSolid class="mr-1" size="lg" />
                    {/if}
                    <P size="xl">{item.Name}</P>
                </div>
            </div>
        </Listgroup>
        {/if}
    </main>
    <div class="fixed bottom-7 right-10">
        <Button pill size="xl" disabled={isExportButtonDisabled()}
                onclick={onExportClick}>Export</Button>
    </div>
</div>