<script lang="ts">
  import "./app.css";
  import {Button, P, Navbar, NavBrand, Input, Label, Helper} from 'flowbite-svelte';
  import { ArrowRightOutline, TabletSolid } from 'flowbite-svelte-icons';
  import {ReadFiles, IsIpValid} from '../wailsjs/go/main/App.js';

  let rmIp: string = $state("10.11.99.1");
  let isRmIpValid: boolean = $state(true);
  $effect(() => {
    if (rmIp) {
      IsIpValid(rmIp).then((result: boolean) => isRmIpValid = result);
    }
  });

</script>

<main class="flex flex-row h-full justify-center items-center">
<div class="flex flex-col flex-wrap content-center items-center justify-between h-1/3 max-h-44">
<Label class="space-y-3 m-5">
  <span>Enter the IP address of your reMarkable tablet:</span>
  <Input placeholder="" size="sm" bind:value={rmIp} color={isRmIpValid ? 'base' : 'red'}>
    <TabletSolid slot="left" class="w-4 h-6" />
  </Input>
</Label>
<Button pill size="xl">
  <ArrowRightOutline class="w-5 h-5 me-2" /> Choose Files
</Button>
</div>

</main>
