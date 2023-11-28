// Right.tsx
import React from 'react';

interface RightProps {
    className: string;
}

const Right: React.FC<RightProps> = ({className}) => {
    return (
        <div className={className}>
            right
        </div>
    );
};

export default Right;
