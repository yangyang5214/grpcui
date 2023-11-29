import React, {useRef, useEffect, useState} from 'react';
import {basicSetup} from 'codemirror';
import {EditorView} from '@codemirror/view';
import {EditorState} from '@codemirror/state';
import {json} from "@codemirror/lang-json"
import {jsonPrettify} from "../utils/json";

// import {oneDark} from '@codemirror/theme-one-dark'; // Import the desired theme


interface JsonViewProps {
    initialDoc: string;
}


export const JsonView: React.FC<JsonViewProps> = ({initialDoc}) => {
    const editorRef = useRef<HTMLDivElement | null>(null);

    const [editor, setEditor] = useState<EditorView | null>(null);

    useEffect(() => {
        if (!editorRef.current) return;
        const formattedDoc = jsonPrettify(initialDoc);

        const state = EditorState.create({
            doc: formattedDoc,
            extensions: [
                basicSetup,
                json(),
                // oneDark,
            ],
        });

        const newEditor = new EditorView({
            state,
            parent: editorRef.current!,
        });

        setEditor(newEditor);

        return () => {
            newEditor.destroy();
        };
    }, [initialDoc]);

    return <div ref={editorRef}></div>;
};
