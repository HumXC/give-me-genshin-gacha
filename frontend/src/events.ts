import { ElMessage } from "element-plus";
import { EventsOn } from "../wailsjs/runtime";
export default function Listen() {
    EventsOn("alert", (m) => ElMessage(m));
}
