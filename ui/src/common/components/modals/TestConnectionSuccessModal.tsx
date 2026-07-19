import { Modal } from "antd"

import { DestinationSuccess } from "@/assets"

const TestConnectionSuccessModal = ({
	open,
	connectionType = "source",
	message,
}: {
	open: boolean
	connectionType?: "source" | "destination" | "catalog"
	// Optional driver warning emitted on success (e.g. a replication slot
	// was created and WAL retention starts now).
	message?: string
}) => {
	const labelMap = {
		source: "Source",
		destination: "Destination",
		catalog: "Catalog",
	} as const
	const label = labelMap[connectionType]
	return (
		<Modal
			open={open}
			footer={null}
			closable={false}
			centered
			width={400}
		>
			<div className="flex flex-col items-center justify-center gap-7 py-6">
				<img src={DestinationSuccess} />
				<div className="flex flex-col items-center">
					<p className="text-xs text-olake-text-tertiary">Successful</p>
					<h2 className="text-lg font-medium">
						{label} test connection is successful
					</h2>
				</div>
				{message && (
					<div className="w-full rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-sm text-amber-700">
						{message}
					</div>
				)}
			</div>
		</Modal>
	)
}

export default TestConnectionSuccessModal
