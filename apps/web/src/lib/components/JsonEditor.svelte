<script lang="ts">
    import CodeMirrorEditor from 'svelte-codemirror-editor';
    import { json } from '@codemirror/lang-json';
    import { EditorView } from '@codemirror/view';
  
    export let value: string;
    export let onChange: (val: string) => void = () => {};
  
    let extensions = [
      json(),
      EditorView.lineWrapping,
      EditorView.theme({
        '&': {
          maxHeight: '400px', // or 600px, 800px, etc.
          border: '1px solid #ccc'
        },
        '.cm-scroller': {
          overflow: 'auto'
        }
      })
    ];
  </script>
  
  <CodeMirrorEditor
    bind:value={value}
    {extensions}
    on:change={(event) => onChange(event.detail.value)}
    class="w-full"
  />