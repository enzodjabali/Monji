<script lang="ts">
    import CodeMirrorEditor from 'svelte-codemirror-editor';
    import { json } from '@codemirror/lang-json';
  
    // The editor’s value (a JSON string)
    export let value: string;
    // Optional change callback (not required, since we use bind:value)
    export let onChange: (val: string) => void = () => {};
  
    // Enable JSON syntax highlighting
    let extensions = [json()];
  </script>
  
  <CodeMirrorEditor
    bind:value={value}
    {extensions}
    on:change={(event) => onChange(event.detail.value)}
    class="border border-gray-300 rounded w-full h-80"
  />