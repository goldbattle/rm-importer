<script lang="ts">
    import { link } from "svelte-spa-router";

    const ipv4Pattern = 
        /^(\d{1,3}\.){3}\d{1,3}$/;
    const ipv6Pattern = 
        /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/;
    let rm_ip: string = $state("10.11.99.1");
    let rm_ip_correct: boolean = $derived(ipv4Pattern.test(rm_ip) || ipv6Pattern.test(rm_ip));
</script>

<div>
    <h1>Welcome</h1>
    <div class="row">
        <div class="col s5 input-field" maxlength="20">
            <div class="prefix"><i class="material-icons">dns</i></div>
            <input id="inp" type="text" bind:value={rm_ip}>
            <label for="inp">IP Address of a Remarkable Tablet</label>
        </div>
        {rm_ip_correct}

        <div class="fixed-action-btn">
            <a class="btn-floating btn-large">
            <i class="large material-icons">mode_edit</i>
            </a>
            <ul>
            <li><a class="btn-floating"><i class="material-icons">insert_chart</i></a></li>
            <li><a class="btn-floating yellow darken-1"><i class="material-icons">format_quote</i></a></li>
            <li><a class="btn-floating green"><i class="material-icons">publish</i></a></li>
            <li><a class="btn-floating blue"><i class="material-icons">attach_file</i></a></li>
            </ul>
        </div>
    </div>
    <a href="/rm-files" use:link>rm-files</a>
</div>