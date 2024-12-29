<script lang="ts">
    import { link } from "svelte-spa-router";
    let rm_ip: string = "10.11.99.1";
    let rm_ip_correct = true;
    const ip_regex = '/^(?>(?>([a-f0-9]{1,4})(?>:(?1)){7}|(?!(?:.*[a-f0-9](?>:|$)){8,})((?1)(?>:(?1)){0,6})?::(?2)?)|(?>(?>(?1)(?>:(?1)){5}:|(?!(?:.*[a-f0-9]:){6,})(?3)?::(?>((?1)(?>:(?1)){0,4}):)?)?(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])(?>\.(?4)){3}))$/iD';

    function on_rm_ip_change() {
        $: rm_ip_correct = (rm_ip.match(ip_regex) != null);
    }
</script>

<div>
    <h1>Welcome</h1>
    <div class="row">
        <div class="col s5 input-field" maxlength="20">
            <div class="prefix"><i class="material-icons">dns</i></div>
            <input id="inp" type="text" value={rm_ip} on:change={on_rm_ip_change}>
            <label for="inp">IP Address of a Remarkable Tablet</label>
        </div>
    </div>
    <a href="/rm-files" use:link>rm-files</a>
</div>