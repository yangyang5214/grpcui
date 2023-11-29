// @ts-ignore
import axios from "axios";

// const apiEndPoint = "http://localhost:8548"
const apiEndPoint = ""


interface AllMethodResponse {
    methods: {
        service: string;
        methods: string[];
    }[];
    addr: string;
}

export async function GetAllMethods() {
    return await axios.get<AllMethodResponse>(apiEndPoint + '/all/methods');
}


interface GetPayloadResponse {
    payload: string
}


export async function GetPayload(method: string) {
    return await axios.get<GetPayloadResponse>(apiEndPoint + '/method/fake_body?method=' + method);
}

export async function Send(method: string, payload: string) {
    return await axios.post<string>(apiEndPoint + '/send', {
        "payload": payload,
        "method": method,
    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    })
}

