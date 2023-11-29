import React, {useEffect, useState} from 'react';
import './App.css';
import {GetAllMethods, GetPayload, Send} from "./Api";
import {JsonView} from "./ui/json-view";


interface Api {
    service: string;
    methods: string[];
}


function App() {

    const [apis, setApis] = useState<Api[]>();
    const [addr, setAddr] = useState("");

    const [payload, setPayload] = useState<string>("");

    const [selectMethod, setSelectMethod] = useState<string>("");

    let respBody: string = ""

    useEffect(() => {
        // Fetch initial data when the component mounts
        const fetchInitialData = async () => {
            try {
                const resp = await GetAllMethods();
                setApis(resp.data.methods);
                setAddr(resp.data.addr);
            } catch (error) {
                console.log(error)
            }
        };

        fetchInitialData().then(r => {
            console.log(r)
        });
    }, []); // Empty dependency array ensures the effect runs once when the component mounts


    const handleMethodClick = async (method: string) => {
        // Handle the click event for the method
        console.log(`Method clicked: ${method}`);
        setSelectMethod(method)
        console.log(`selectMethod is: `, selectMethod)
        await getPayload()
    };

    async function sendHttp() {
        console.log("do send ...")
        const resp = await Send(selectMethod, payload)
        respBody = resp.data
    }


    async function getPayload() {
        if ("" === selectMethod) {
            return
        }
        const resp = await GetPayload(selectMethod)
        setPayload(resp.data.payload)
        console.log(`payload is`, payload)
    }


    const apiItems = apis?.map(api => (
        <div key={api.service} className="api-container">
            {api.service}
            <div className="method-container">
                {api.methods.map(method => (
                    <div key={method} className="method" onClick={() => handleMethodClick(method)}>
                        {method}
                    </div>
                ))}
            </div>
        </div>
    ));

    return (
        <div className="App">
            <div className={"Left"}>
                {apiItems}
            </div>

            <div className={"Send"}>
                <div className="endpoint">
                    <span className="methodName">{addr} {selectMethod} </span>
                    <button onClick={sendHttp}> Send</button>
                </div>
                <div className="payload">
                    <JsonView initialDoc={payload}/>
                </div>
            </div>

            <div className={"Right"}>
                <div>
                    {respBody}
                </div>
            </div>
        </div>
    );
}

export default App;
