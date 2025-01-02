import { writable, type Writable } from "svelte/store";

const filesToExport: Writable<string[]> = writable([]);

function createCheckboxSignal(): [(key: string) => boolean, 
                                  (key: string, value: boolean) => void,
                                  () => boolean,
                                  () => void] {
    let checkbox: { [key: string]: boolean; } = $state({});
    let checkedCount: number = $state(0);

    return [
            (key: string) => checkbox[key], 
            (key: string, value: boolean) => {
                if (value === true) {
                    checkedCount += 1;
                } else {
                    checkedCount -= 1;
                }
                checkbox[key] = value;
            },
            () => {
                return checkedCount === 0;
            },
            () => {
                let result: string[] = [];
                for (let key in checkbox) {
                    if (checkbox[key]) {
                        result.push(key);
                    }
                }
                filesToExport.set(result);
            }
           ];
}

export const [isChecked, setChecked, isExportButtonDisabled, storeCheckedFiles] = createCheckboxSignal();
export {filesToExport};