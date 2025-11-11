const owl = {
    title: 'Excelent owl',
    src: 'https://content.codecademy.com/courses/React/react_photo-owl.jpg',
};

export default function Coruja() {
    return (
        <>
            <h1>{owl.title}</h1>
            <img
                src={owl.src}
                alt={owl.title}
            />
        </>
    );
}