<script lang="ts">
    import { Button, Checkbox, Heading, Listgroup, Navbar, P, Toolbar, ToolbarButton } from "flowbite-svelte";
    import { ArrowLeftOutline, ArrowUpOutline, CheckOutline, FileLinesOutline, FileLinesSolid, FileOutline, FolderArrowRightSolid, FolderOutline, FolderSolid, ImageOutline, MapPinAltSolid, PaperClipOutline } from "flowbite-svelte-icons";
    import { push, location, link } from "svelte-spa-router";
    import { GetTabletFolder } from "../../wailsjs/go/main/App";

    let id = $state("");
    let path: string[] = $state([]);
    let links: DocInfo[] = $state([]);
    let checkboxes = $state({});
    let checked_checkboxes = $state(0);

    let export_button_disabled = $derived.by(() => {
        return checked_checkboxes == 0;
    });

    $effect(() => {
        GetTabletFolder(id).then((result) => {
            links = result;
        });
    });

    const onBack = () => {
        if (path.length > 0) {
            id = path[path.length - 1];
            path.pop();
        } else {
            id = '';
        }
    }

    const onCheckboxClick = (item: DocInfo, checked: boolean) => {
        checkboxes[item] = checked;
        if (checked) {
            checked_checkboxes += 1;
        } else {
            checked_checkboxes -= 1;
        }
    };
</script>

<div style="height: fit-content;">
    <Navbar color="blue" class="sticky top-0">
        <div class={path.length == 0 ? "invisible": ""}>
            <ToolbarButton color="blue" name="Back" onclick={onBack}> <ArrowUpOutline class="w-7 h-7" /></ToolbarButton>
        </div>
        <h1 class="font-bold">Choose items to export</h1>
        <div>
            <ToolbarButton class="invisible"><ArrowUpOutline class="w-7 h-7" /></ToolbarButton>
        </div>
    </Navbar>
    <main class="pl-10 pr-10 pt-3 pb-3">
        {#if links.length > 0}
        <Listgroup items={links} let:item active={false}>
            <div class="flex flex-row justify-start items-center">
                <Checkbox checked={checkboxes[item.Id]} 
                          on:click={(e) => onCheckboxClick(item.Id, e.target.checked)} 
                class="mr-2" />
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
        <Button pill size="xl" disabled={export_button_disabled}>Export</Button>
    </div>
</div>