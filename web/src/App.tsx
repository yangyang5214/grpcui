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

    const [visibleApis, setVisibleApis] = useState<Array<String>>([]);

    useEffect(() => {
        console.log(`useEffect method call...`)
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
        if (selectMethod === "") {
            alert("Need select method")
            return
        }
        clear()
        await Send(selectMethod, payload).then(resp => {
            setRespBody(resp.data)
            setRespCode(resp.status + " " + resp.statusText)
        }).catch(err => {
            let resp = err.response
            if (resp === undefined) {
                alert("server internal error")
                return
            }
            console.log(resp)
            setRespBody(resp.data)
            setRespCode(resp.status + " " + resp.statusText)
        })
    }


    async function getPayload(method: string) {
        const resp = await GetPayload(method)
        setPayload(resp.data.payload)
    }


    function handleApiClick(service: string) {
        if (visibleApis.includes(service)) {
            // If it is, remove it
            setVisibleApis((prevArray) => prevArray.filter((existingItem) => existingItem !== service));
        } else {
            // If it's not, append it
            setVisibleApis((prevArray) => [...prevArray, service]);
        }
    }

    const apiItems = apis?.map(api => (
        <div>
            <div key={api.service} className="api-container" onClick={() => handleApiClick(api.service)}>
                {api.service}
            </div>
            <div className={`method-container ${visibleApis.includes(api.service) ? '' : 'hide-display'}`}>
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
                    <div className={"statusCode"}>
                        {respCode}
                    </div>

                    <div>
                        {respBody !== "" && <JsonView initialDoc={respBody}/>}
                    </div>
                </div>
            </div>
        </div>
    )
        ;
}

export default App;
