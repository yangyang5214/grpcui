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

    const [respBody, setRespBody] = useState<string>("");
    const [respCode, setRespCode] = useState<String>("");


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

    function clear() {
        setRespBody("")
        setRespCode("")
    }

    const handleMethodClick = async (method: string) => {
        // Handle the click event for the method
        console.log(`Method clicked: ${method}`);

        clear()

        console.log(`selectMethod is: `, method)

        await getPayload(method)
        setSelectMethod(method)
    };

    async function sendHttp() {
        clear()
        await Send(selectMethod, payload).then(resp => {
            setRespBody(resp.data)
            setRespCode(resp.status + " " + resp.statusText)
        }).catch(err => {
            let resp = err.response
            setRespBody(resp.data)
            setRespCode(resp.status + " " + resp.statusText)
        })
    }


    async function getPayload(method: string) {
        const resp = await GetPayload(method)
        setPayload(resp.data.payload)
        console.log(`payload is`, payload)
    }


    const apiItems = apis?.map(api => (
        <div key={api.service} className="api-container">
            {api.service}
            <div className="method-container">
                {api.methods.map(method => {
                    const parts = method.split('.');
                    const methodName = parts[parts.length - 1];
                    return <div key={method} className="method" onClick={() => handleMethodClick(method)}>
                        {methodName}
                    </div>
                })}
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
                    <div className="addr">{addr} </div>
                    <div className="methodName"> {selectMethod} </div>
                    <button className="send" onClick={sendHttp}> Send</button>
                </div>
                <div className="payload">
                    <JsonView initialDoc={payload}/>
                </div>
            </div>

            <div className={"Right"}>
                <div>
                    <div>
                        {respCode}
                    </div>

                    <div>
                        {respBody}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default App;
