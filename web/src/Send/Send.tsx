import React from 'react';


interface SendProps {
    className: string;
}

const Send: React.FC<SendProps> = ({className}) => {
    return (
        <div className={className}>
            send
        </div>
    );
}

export default Send;
