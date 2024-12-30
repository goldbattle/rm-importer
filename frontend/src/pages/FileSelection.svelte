<script lang="ts">
    import { Button, Checkbox, Heading, Listgroup, Navbar, P, Toolbar, ToolbarButton } from "flowbite-svelte";
    import { ArrowLeftOutline, ArrowUpOutline, CheckOutline, FileLinesOutline, FileLinesSolid, FileOutline, FolderArrowRightSolid, FolderOutline, ImageOutline, MapPinAltSolid, PaperClipOutline } from "flowbite-svelte-icons";
    import { push } from "svelte-spa-router";

    export let params = {}

    let links = [
        { name: 'Accordions0', href: '#/files/componentsaccordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges2', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
        { name: 'Accordions', href: '/docs/components/accordion'},
        { name: 'Alerts', href: '/docs/components/alert' },
        { name: 'Badges', href: '/docs/components/badge'},
        { name: 'Breadcrumbs', href: '/docs/components/breadcrumb', attrs: {target: '_blank'} },
    ];
</script>

<div>
    <div class="sticky top-0">
        <Navbar color="blue" class="sticky top-0">
            <div>
                <ToolbarButton color="blue" name="Attach file" onclick={() => push('/')}><ArrowLeftOutline class="w-7 h-7" /></ToolbarButton>
            </div>
            <h1 class="font-bold">Choose items to export</h1>
            <div>
                <ToolbarButton color="blue" name="Upload image"><ArrowUpOutline class="w-7 h-7" /></ToolbarButton>
            </div>
        </Navbar>
    </div>
    <main class="pl-10 pr-10 pt-3">
        <Listgroup items={links} let:item>
            <div class="flex flex-row justify-start items-center">
                <Checkbox class="mr-2"></Checkbox>
                <div class="flex flex-row justify-start items-center hover:bg-gray-200 w-full">
                    <FileLinesSolid class="mr-1" size="lg" />
                    <P size="xl">{item.name}</P>
                </div>
            </div>
        </Listgroup>
    </main>
    <div hidden={false} class="fixed bottom-7 right-10">
        <Button pill size="xl">Export</Button>
    </div>
</div>