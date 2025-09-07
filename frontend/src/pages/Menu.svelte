<script lang="ts">
  import { Alert, Button, P, Input, Label, Spinner, Footer, A, Select, Checkbox} from 'flowbite-svelte';
  import { ArrowRightOutline, InfoCircleSolid, TabletSolid, CloseOutline, ServerSolid, UserSolid } from 'flowbite-svelte-icons';
  import { ReadDocs, IsIpValid, GetAppVersion, ConnectSSH, ConnectSSHForUploads, GetSafeMode, SetSafeMode, SetHybridMode, TestSSHConnection } from '../../wailsjs/go/main/App.js';
  import { push } from 'svelte-spa-router';
  import { BrowserOpenURL } from '../../wailsjs/runtime/runtime.js';

  const source = "https://github.com/chopikus/rm-importer";
  let version = $state("");

  GetAppVersion().then((v: string) => {
    version = v;
  });
  
  // Connection mode
  let useSSH: boolean = $state(true);
  let hybridMode: boolean = $state(false);
  
  // HTTP connection (legacy)
  let rmIp: string = $state("10.11.99.1");
  let isRmIpValid: boolean = $state(true);
  
  // SSH connection
  let sshHost: string = $state("10.11.99.1");
  let sshUsername: string = $state("root");
  let sshPassword: string = $state("");
  
  let loading: boolean = $state(false);
  let show_error: boolean = $state(false);
  let error_message: string = $state("");
  let safe_mode: boolean = $state(true);

  // Load safe mode setting
  GetSafeMode().then((mode: boolean) => {
    safe_mode = mode;
  });

  $effect(() => {
    if (rmIp) {
      IsIpValid(rmIp).then((result: boolean) => isRmIpValid = result);
    }
  });

  function onNext() {
    if (loading) {
      return;
    }
    show_error = false;
    loading = true;
    
    if (useSSH) {
      if (hybridMode) {
        // Hybrid mode: HTTP for listing, SSH for uploads
        // First connect via HTTP for file listing
        ReadDocs(sshHost)
          .then((_: any) => {
            // Then establish SSH connection for uploads only
            return ConnectSSHForUploads(sshHost, sshUsername, sshPassword);
          })
          .then((_: any) => push('/files'))
          .catch((err: Error) => {
            console.log("Couldn't connect in hybrid mode! Make sure both HTTP and SSH are working.");
            error_message = err.toString()
            show_error = true
          })
          .finally(() => {
            loading = false;
          });
      } else {
        // Pure SSH connection
        ConnectSSH(sshHost, sshUsername, sshPassword)
          .then((_: any) => push('/files'))
          .catch((err: Error) => {
            console.log("Couldn't connect to reMarkable tablet via SSH! Make sure the credentials are correct.");
            error_message = err.toString()
            show_error = true
          })
          .finally(() => {
            loading = false;
          });
      }
    } else {
      // HTTP connection (legacy)
      ReadDocs(rmIp)
        .then((_: any) => push('/files'))
        .catch((err: Error) => {
          console.log("Couldn't connect to reMarkable tablet! Make sure the IP address is correct.");
          error_message = err.toString()
          show_error = true
        })
        .finally(() => {
          loading = false;
        });
    }
  }

  function onSafeModeToggle() {
    safe_mode = !safe_mode;
    SetSafeMode(safe_mode);
  }

  function onHybridModeToggle() {
    hybridMode = !hybridMode;
    SetHybridMode(hybridMode);
  }

  function onTestSSH() {
    if (loading) {
      return;
    }
    show_error = false;
    loading = true;
    
    TestSSHConnection(sshHost, sshUsername, sshPassword)
      .then((_: any) => {
        alert("SSH connection test successful!");
      })
      .catch((err: Error) => {
        console.log("SSH connection test failed:", err);
        error_message = err.toString()
        show_error = true
      })
      .finally(() => {
        loading = false;
      });
  }

  const onSourceClick = () => {
    BrowserOpenURL(source)
  };

</script>

{#key show_error}
{#if show_error}
<Alert color="red" dismissable={true}
  class="border-t-4 top-0 absolute left-1/2 transform -translate-x-1/2 w-full">
  <div class="flex items-center gap-1">
    <InfoCircleSolid class="w-5 h-5" />
    <span class="text-lg font-medium">Couldn't connect to reMarkable!</span>
  </div>
  <p class="mt-2 mb-4 text-sm">{error_message}</p>
</Alert>
{/if}
{/key}

<main class="flex flex-col h-full justify-center items-center content-center">
  <P size="3xl">rm-importer</P>
  
  <!-- Settings Panel -->
  <div class="w-96 mb-6 p-4 bg-gray-50 rounded-lg border">
    <h3 class="text-lg font-semibold mb-4 text-gray-800">Settings</h3>
    
    <!-- Connection Mode Toggle -->
    <div class="flex items-center justify-between mb-3">
      <span class="text-sm font-medium text-gray-700">Use SSH Mode?:</span>
      <Checkbox bind:checked={useSSH} color="blue">
        <span slot="label">{useSSH ? 'SSH' : 'HTTP'}</span>
      </Checkbox>
    </div>

    <!-- Hybrid Mode Toggle (only visible when SSH is selected) -->
    {#if useSSH}
    <div class="flex items-center justify-between mb-3">
      <div class="flex flex-col">
        <span class="text-sm font-medium text-gray-700">Hybrid Mode:</span>
        <span class="text-xs text-gray-500">HTTP for listing, SSH for uploads</span>
      </div>
      <Checkbox checked={hybridMode} color="purple" on:change={onHybridModeToggle}>
        <span slot="label">{hybridMode ? 'Enabled' : 'Disabled'}</span>
      </Checkbox>
    </div>
    {/if}

    <!-- Safe Mode Toggle -->
    <div class="flex items-center justify-between">
      <div class="flex flex-col">
        <span class="text-sm font-medium text-gray-700">Safe Mode:</span>
        <span class="text-xs text-gray-500">Prevents accidental uploads</span>
      </div>
      <Checkbox checked={safe_mode} color="green" on:change={onSafeModeToggle}>
        <span slot="label">{safe_mode ? 'Enabled' : 'Disabled'}</span>
      </Checkbox>
    </div>
  </div>

  <div class="flex flex-col flex-wrap content-center items-center justify-between min-h-64 w-96">
    {#if useSSH}
      <!-- SSH Connection Form -->
      <div class="space-y-4 w-full">
        <Label class="space-y-2">
          <span>SSH Host (IP Address):</span>
          <Input placeholder="10.11.99.1" size="sm" bind:value={sshHost}>
            <ServerSolid slot="left" class="w-4 h-4" />
          </Input>
        </Label>

        <Label class="space-y-2">
          <span>Username:</span>
          <Input placeholder="root" size="sm" bind:value={sshUsername}>
            <UserSolid slot="left" class="w-4 h-4" />
          </Input>
        </Label>

        <Label class="space-y-2">
          <span>Password:</span>
          <Input type="password" placeholder="Enter SSH password" size="sm" bind:value={sshPassword}>
          </Input>
        </Label>
      </div>
    {:else}
      <!-- HTTP Connection Form (Legacy) -->
      <Label class="space-y-3 m-5">
        <span>Enter the IP address of your reMarkable tablet:</span>
        <Input placeholder="" size="sm" bind:value={rmIp} color={isRmIpValid ? 'base' : 'red'}>
          <TabletSolid slot="left" class="w-4 h-6" />
        </Input>
      </Label>
    {/if}

    <div class="flex gap-3 mt-6">
      {#if useSSH}
        <Button size="lg" color="gray" on:click={onTestSSH} disabled={loading}>
          {#if !loading}
            Test SSH
          {:else}
            <Spinner />
          {/if}
        </Button>
      {/if}
      
      <Button pill size="xl" on:click={onNext} disabled={loading}>
        {#if !loading}
          <ArrowRightOutline class="w-5 h-5 me-2" /> Connect
        {:else}
          <Spinner />
        {/if}
      </Button>
    </div>
  </div>
</main>
<Footer class="absolute bottom-0 left-0 right-0 flex justify-center items-center py-4">
  <P size="sm">Version {version}. Source available at <A onclick={onSourceClick}>{source}</A></P>
</Footer>
