<script lang="ts">
    import { Button, Checkbox, Listgroup, Navbar, P, ToolbarButton } from "flowbite-svelte";
    import { ArrowUpOutline, FileLinesSolid, FolderSolid } from "flowbite-svelte-icons";
    import { GetTabletFolder, GetTabletFolderSelection, OnItemSelect, GetCheckedFilesCount } from "../../wailsjs/go/main/App";
    import { push } from "svelte-spa-router";
    import { backend } from "../../wailsjs/go/models";
    type DocInfo = backend.DocInfo;
    
    let id = $state("");
    let path: string[] = $state([]);
    let items: DocInfo[] = $state([]);

    // Checked stores checkbox value for every element of the folder
    let checked: {[key: string]: number} = $state({});
    // defined in go code
    const UNSELECTED = 0, INDETERMINATE = 1, SELECTED = 2;

    let export_disabled = $state(true);

    // onIdUpdate
    $effect(() => {
        GetTabletFolder(id).then((result) => {
            items = result;
        });
        GetTabletFolderSelection(id).then((result) => {
            checked = {}
            for (const item of result) {
                checked[item.Id] = item.Status
            }
        });
    });

    const checkUpdate = (item: DocInfo, value: boolean | undefined) => {
        let select;
        if (value) {
            checked[item.Id] = SELECTED;
            select = true;
        } else {
            checked[item.Id] = UNSELECTED;
            select = false;
        }
        OnItemSelect(item.Id, select)
        .then(() => GetCheckedFilesCount().then((count: number) => {
                export_disabled = (count === 0);
            }));
    };

    const onBack = () => {
        if (path.length > 0) {
            id = path[path.length - 1];
            path.pop();
        } else {
            id = '';
        }
    };

    const onItemClick = (item: DocInfo) => {
        if (item.IsFolder) {
            path.push(id);
            id = item.Id;
        }
    };

    const onExportClick = () => {
        //storeCheckedFiles();
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
                
                {#key [id, checked[""]]}
                <Checkbox bind:checked={() => checked[item.Id] === SELECTED, (v) => checkUpdate(item, v)}
                          indeterminate={checked[item.Id] === INDETERMINATE}
                          class="mr-2" />
                {/key}

                <div class="flex flex-row justify-start items-center w-full hover:bg-gray-100"
                     onclick={() => onItemClick(item)}>
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
        <Button pill size="xl" disabled={export_disabled}
                onclick={onExportClick}>Export</Button>
    </div>
</div>