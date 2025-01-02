<script lang="ts">
  import "./app.css";
  import Menu from "./pages/Menu.svelte";
  import FileSelection from "./pages/FileSelection.svelte";
  import Router from "svelte-spa-router";
    import ExportConfirmation from "./pages/ExportConfirmation.svelte";

  const routes = {
    '/': Menu,
    '/files': FileSelection,
    '/export-confirmation': ExportConfirmation,
    '*': Menu,
  };
</script>

<Router {routes}/>