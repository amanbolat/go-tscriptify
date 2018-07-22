/* Do not change, this code is generated from Golang structs */


export interface Expense {
    payment_on_deliver: Money;
    //[Expense:]


    //[end]
}


export interface WarehouseEntry {
    id: string;
    shipment_number: string;
    date_of_entry: Date;
    source: string;
    track_code: string;
    box_qty: number;
    pcs_qty: number;
    product_name: string;
    //[WarehouseEntry:]


    //[end]
}
export interface Decimal {
    //[Decimal:]


    //[end]
}
export interface UnitLoad {
    sequence: number;
    quantity: number;
    product_name: string;
    weight: Decimal;
    length: number;
    height: number;
    width: number;
    //[UnitLoad:]


    //[end]
}

export interface Recipient {
    name: string;
    phone_numbers: string[];
    address: Address;
    //[Recipient:]


    //[end]
}
export interface Address {
    country: string;
    state: string;
    city: string;
    street: string;
    houseNumber: string;
    //[Address:]


    //[end]
}
export interface Sender {
    name: string;
    phone_numbers: string[];
    address: Address;
    //[Sender:]


    //[end]
}









export interface Money {
    //[Money:]


    //[end]
}


export enum PackageMethod {
    PackageNone = 'package_none',
    PackageBag = 'package_bag',
    PackageStandard = 'package_standard',
    PackageCarton = 'package_carton',
    PackageFoam = 'package_foam',
    PackageCartonFoam = 'package_carton_foam',
    PackageWoodenCrate = 'package_wooden_crate',
    PackageWoodenCrateFoam = 'package_wooden_crate_foam',
    PackageWoodenBox = 'package_wooden_box',
    PackageWoodenBoxFoam = 'package_wooden_box_foam',
    //[PackageMethod:]


    //[end]
}
export enum PaymentStatus {
    Unpaid = 'unpaid',
    OnCredit = 'on_credit',
    Paid = 'paid',
    //[PaymentStatus:]


    //[end]
}
export enum PaymentMethod {
    OnDelivery = 'on_delivery',
    PersonalTransfer = 'personal_transfer',
    LegalEntityTransfer = 'legal_entity_transfer',
    Other = 'other',
    //[PaymentMethod:]


    //[end]
}
export enum TransportMethod {
    Air = 'air',
    Auto = 'auto',
    Train = 'train',
    Express = 'express',
    Sea = 'sea',
    Local = 'local',
    //[TransportMethod:]


    //[end]
}
export enum ShipmentStatus {
    Planning = 'planning',
    Preparation = 'preparation',
    Packed = 'packed',
    SentOut = 'sent_out',
    CustomsClearance = 'customs_clearance',
    OnTheWayToTp = 'on_the_way_to_tp',
    DeliveredToTp = 'delivered_to_tp',
    DeliveredToRecipient = 'delivered_to_recipient',
    InvalidStatus = 'invalid_status',
    //[ShipmentStatus:]


    //[end]
}
export enum ShipmentType {
    CommonShipment = 'common_shipment',
    ConsolidationShipment = 'consolidation_shipment',
    //[ShipmentType:]


    //[end]
}
export interface Shipment {
    id: string;
    code: string;
    type: ShipmentType;
    customer_code: string;
    packages_qty: number;
    pieces_qty: number;
    current_status: ShipmentStatus;
    transport_method: TransportMethod;
    payment_method: PaymentMethod;
    payment_status: PaymentStatus;
    package_method: PackageMethod;
    departure_warehouse: string;
    arrival_warehouse: string;
    transfer_point_warehouse: string;
    date_created: Date;
    date_modified: Date;
    value: Money;
    package_fee: Money;
    transfer_fee: Money;
    freight_rate: Money;
    prepayment: Money;
    discount: Money;
    commission_fee: Money;
    total_for_payment: Money;
    insurance_rate: number;
    insurance_fee: Money;
    payment_on_deliver: Money;
    sender: Sender;
    recipient: Recipient;
    unit_loads: UnitLoad[];
    entries: WarehouseEntry[];
    consolidation: Shipment[];
    expense: Expense;
    image_urls: string[];
    //[Shipment:]


    //[end]
}
export interface Address {
    city: string;
    number: number;
    country: string;
    //[Address:]


    //[end]
}
export interface PersonalInfo {
    hobby: string[];
    pet_name: string;
    //[PersonalInfo:]


    //[end]
}
export interface Person {
    name: string;
    personal_info: PersonalInfo;
    nicknames: string[];
    addresses: Address[];
    shipments: Shipment[];
    //[Person:]


    //[end]
}