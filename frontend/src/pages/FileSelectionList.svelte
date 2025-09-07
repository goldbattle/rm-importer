<script lang="ts">
    import { Listgroup, Checkbox, P, Button } from "flowbite-svelte";
    import { FolderSolid, FileLinesSolid, ArrowUpOutline, InfoCircleSolid } from "flowbite-svelte-icons";
    import { backend } from "../../wailsjs/go/models";
    import { FileDialog, UploadFileSSH, GetSafeMode, IsSSHMode } from "../../wailsjs/go/main/App.js";
    type DocInfo = backend.DocInfo;

    let {items, isItemChecked, isItemIndeterminate, itemCheckUpdate, onItemClick, folderId, addItemToList} = $props();
    
    let safe_mode: boolean = $state(true);
    let ssh_mode: boolean = $state(false);
    
    // Notification system - array to handle multiple notifications
    let notifications: Array<{id: string, message: string, type: 'success' | 'error' | 'info'}> = $state([]);
    
    // Upload state management
    let isUploading: boolean = $state(false);
    
    // Load safe mode and SSH mode settings
    GetSafeMode().then((mode: boolean) => {
        safe_mode = mode;
    });
    
    IsSSHMode().then((mode: boolean) => {
        ssh_mode = mode;
        console.log("SSH mode detected:", mode);
    });

    function showNotification(message: string, type: 'success' | 'error' | 'info' = 'info') {
        const id = `notification_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        
        notifications = [...notifications, {
            id: id,
            message: message,
            type: type
        }];
        
        // Auto-hide after 8 seconds for success, 10 seconds for errors, 5 seconds for info
        const timeout = type === 'success' ? 8000 : type === 'error' ? 10000 : 5000;
        setTimeout(() => {
            removeNotification(id);
        }, timeout);
    }

    function removeNotification(id: string) {
        notifications = notifications.filter(n => n.id !== id);
    }

    async function copyUUID(uuid: string, itemName: string) {
        try {
            await navigator.clipboard.writeText(uuid);
            showNotification(`UUID copied to clipboard: ${itemName}`, 'success');
        } catch (error) {
            console.error('Failed to copy UUID:', error);
            showNotification(`Failed to copy UUID: ${error}`, 'error');
        }
    }

    async function onUploadClick() {
        if (isUploading) return; // Prevent multiple uploads
        
        // Check current safe mode status
        const currentSafeMode = await GetSafeMode();
        if (currentSafeMode) {
            showNotification("Upload blocked: Safe mode is enabled. Please disable safe mode in the main menu to upload files.", 'error');
            return;
        }
        
        console.log("Upload button clicked for folder:", folderId);
        
        const filePath = await FileDialog();
        console.log("FileDialog returned:", filePath);
        
        if (filePath && filePath.trim() !== "") {
            try {
                isUploading = true;
                
                // Extract filename from path, handling both Windows and Unix paths
                const pathParts = filePath.split(/[\\/]/);
                const fileName = pathParts[pathParts.length - 1] || "";
                
                console.log("Extracted filename:", fileName);
                
                if (!fileName) {
                    showNotification("Invalid file path selected", 'error');
                    return;
                }
                
                // Show upload progress
                showNotification(`Uploading "${fileName}"... Please wait.`, 'info');
                console.log("Starting upload:", filePath, "to folder:", folderId);
                
                await UploadFileSSH(filePath, fileName, folderId);
                console.log("Upload completed successfully");
                
                // Create a new DocInfo object for the uploaded file
                const visibleName = fileName.replace(/\.pdf$/i, ''); // Remove .pdf extension for display
                const newFile: DocInfo = {
                    Id: crypto.randomUUID(), // Generate a temporary UUID for the UI
                    Name: visibleName, // Use the name without extension
                    IsFolder: false,
                    ParentId: folderId,
                    Type: "DocumentType",
                    Version: 1,
                    Modified: new Date().toISOString(),
                    Size: 0 // We don't know the actual size yet
                };
                
                // Add the new file to the list
                addItemToList(newFile);
                
                showNotification(`File "${fileName}" uploaded successfully! It may take a moment to appear on the device after xochitl restarts.`, 'success');
            } catch (error) {
                console.error("Upload failed:", error);
                showNotification(`Upload failed: ${error}`, 'error');
            } finally {
                isUploading = false;
            }
        } else if (filePath === "") {
            // User cancelled the dialog
            console.log("User cancelled file dialog");
            return;
        } else {
            console.log("No file selected or empty path");
            showNotification("No file selected", 'error');
        }
    }

</script>

<!-- Upload Button at the top -->
{#if ssh_mode}
    <div class="mb-4 flex justify-center">
        <Button 
            color="blue" 
            size="lg" 
            disabled={isUploading || safe_mode}
            onclick={onUploadClick}
            class="px-6 py-3">
            {#if isUploading}
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Uploading...
            {:else}
                <ArrowUpOutline class="w-5 h-5 mr-2" />
                Upload File to {folderId ? 'Current Folder' : 'Root Directory'}
            {/if}
        </Button>
    </div>
    
    {#if safe_mode}
        <div class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
            <p class="text-sm text-yellow-800">
                <strong>Safe Mode Enabled:</strong> Upload is disabled. Please disable safe mode in the main menu to upload files.
            </p>
        </div>
    {/if}
{/if}

{#if items.length > 0}
        <Listgroup {items} let:item active={false}>
            <div class="flex flex-row justify-start items-center">
                <Checkbox bind:checked={() => isItemChecked(item.Id), (v) => itemCheckUpdate(item.Id, v)}
                          indeterminate={isItemIndeterminate(item.Id)}
                          class="mr-4 w-4 h-4" />

                <div class="flex flex-row justify-start items-center w-full hover:bg-gray-100 cursor-pointer"
                     role="button" tabindex="0"
                     onclick={() => onItemClick(item)}
                     onkeydown={(e) => e.key === 'Enter' && onItemClick(item)}>
                    {#if item.IsFolder}
                        <FolderSolid class="mr-0.5" size="lg" />
                    {:else}
                        <FileLinesSolid class="mr-0.5" size="lg" />
                    {/if}
                    <P size="xl">{item.Name}</P>
                </div>
                
                <!-- Copy UUID Button -->
                <button 
                    class="ml-auto mr-2 p-1 text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors duration-200"
                    onclick={(e) => {
                        e.stopPropagation();
                        copyUUID(item.Id, item.Name);
                    }}
                    title="Copy UUID to clipboard"
                    aria-label="Copy UUID to clipboard">
                    <InfoCircleSolid class="w-4 h-4" />
                </button>
            </div>
        </Listgroup>
{/if}

<!-- Notification System -->
{#if notifications.length > 0}
    <div class="fixed top-0 left-0 right-0 z-50 p-4 space-y-2">
        {#each notifications as notification (notification.id)}
            <div class="mx-auto max-w-4xl">
                <div class="p-4 rounded-lg shadow-lg border-l-4 {
                    notification.type === 'success' ? 'bg-green-50 border-green-400 text-green-800' :
                    notification.type === 'error' ? 'bg-red-50 border-red-400 text-red-800' :
                    'bg-blue-50 border-blue-400 text-blue-800'
                }">
                    <div class="flex items-center">
                        <div class="flex-shrink-0">
                            {#if notification.type === 'success'}
                                <svg class="h-6 w-6 text-green-400" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                                </svg>
                            {:else if notification.type === 'error'}
                                <svg class="h-6 w-6 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                                </svg>
                            {:else}
                                <svg class="h-6 w-6 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                                </svg>
                            {/if}
                        </div>
                        <div class="ml-3 flex-1">
                            <p class="text-base font-medium">{notification.message}</p>
                        </div>
                        <div class="ml-auto pl-3">
                            <button onclick={() => removeNotification(notification.id)} 
                                    aria-label="Close notification"
                                    class="inline-flex rounded-md p-1.5 {
                                notification.type === 'success' ? 'text-green-500 hover:bg-green-100' :
                                notification.type === 'error' ? 'text-red-500 hover:bg-red-100' :
                                'text-blue-500 hover:bg-blue-100'
                            }">
                                <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        {/each}
    </div>
{/if}