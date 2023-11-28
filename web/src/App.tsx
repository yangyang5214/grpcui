import React, {useEffect, useState} from 'react';
import './App.css';
import axios from "axios";


interface ApiResponse {
    methods: {
        service: string;
        methods: string[];
    }[];
    addr: string;
}

interface Api {
    service: string;
    methods: string[];
}


function App() {

    const [apis, setApis] = useState<Api[]>();
    const [addr, setAddr] = useState("");

    const [selectMethod, setSelectMethod] = useState("");

    let respBody = ""

    useEffect(() => {
        // Fetch initial data when the component mounts
        const fetchInitialData = async () => {
            try {
                const resp = await axios.get<ApiResponse>('/all/methods');
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


    const handleMethodClick = (method: string) => {
        // Handle the click event for the method
        console.log(`Method clicked: ${method}`);
        setSelectMethod(method)
    };

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


    function sendHttp() {
        console.log("do send ...")
        respBody = ""
    }

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
