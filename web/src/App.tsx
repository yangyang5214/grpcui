import React from 'react';
import './App.css';
import Left from "./left/Left";
import Send from "./Send/Send";
import Right from "./right/Right";

function App() {
    return (
        <div className="App">
            <Left className={"Left"}/>
            <Send className={"Send"}/>
            <Right className={"Right"}/>
        </div>
    );
}

export default App;
